package main

import (
	"SamporDoc/backend/config"
	"SamporDoc/backend/excel"
	"SamporDoc/backend/infra/model"
	"SamporDoc/backend/infra/repository"
	"SamporDoc/backend/infra/seed"
	"SamporDoc/backend/log"
	"SamporDoc/backend/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type App struct {
	ctx  context.Context
	repo *repository.Repo
}

func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// use wails' context
	a.ctx = ctx
	// connect to the database
	dbFilePath, err := config.GetAppFilePath(config.DBFileName)
	if err != nil {
		fmt.Println(err)
		panic("Error getting database file path")
	}
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}
	// migrate
	db.AutoMigrate(&model.Customer{}, &model.Shop{}, &model.Log{})

	// create repo instance
	a.repo = repository.NewRepo(ctx, db)

	// seeding
	customers, err := a.repo.GetAllCustomers()
	if err != nil {
		fmt.Println("Error seeding database -> repo.GetAllCustomers")
		panic(err)
	}
	if len(customers) == 0 {
		seed.SeedCustomers(a.repo)
	}
	shops, err := a.repo.GetAllShops()
	if err != nil {
		fmt.Println("Error seeding database -> repo.GetAllShops")
		panic(err)
	}
	if len(shops) == 0 {
		seed.SeedShops(a.repo)
	} else {
		seed.SeedUpdateShops(a.repo)
	}
}

// Binding
func (a *App) GetAllShops() ([]model.Shop, error) {
	return a.repo.GetAllShops()
}

func (a *App) UpdateShopBySlug(shop model.Shop) (model.Shop, error) {
	return a.repo.UpdateShopBySlug(&shop)
}

func (a *App) GetAllCustomers() ([]model.Customer, error) {
	return a.repo.GetAllCustomers()
}

func (a *App) GetNextControlNumber(controlFilePath string) (int, error) {
	return excel.GetNextControlNumber(controlFilePath)
}

func (a *App) GetControlData(controlFilePath string, NO int) (*excel.ControlDataWithRowCoordinate, error) {
	return excel.GetControlData(controlFilePath, NO)
}

func (a *App) OpenDirectoryDialog() (string, error) {
	return wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{})
}

func (a *App) OpenExcelFileDialog() (string, error) {
	return wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Filters: []wailsRuntime.FileFilter{{
			DisplayName: "Excel",
			Pattern:     "*.xlsx",
		}},
	})
}

type ReceiptData struct {
	ReceiptNO    string
	ReceiptDate  *time.Time
	CustomerName string
	Address      *string
	Detail       string
	DeliveryNO   *string
	DeliveryDate *time.Time
	Amount       float64
}

func (r ReceiptData) toExcelKeyValue() map[string]string {
	var receiptDate, address, deliveryNO, deliveryDate string
	if r.ReceiptDate == nil {
		receiptDate = strings.Repeat(".", 40)
	} else {
		receiptDate = utils.GetFullThaiDate(*r.ReceiptDate)
	}
	if r.Address != nil {
		address = *r.Address
	}
	if r.DeliveryNO != nil {
		deliveryNO = *r.DeliveryNO
	}
	if r.DeliveryDate != nil {
		deliveryDate = utils.GetShortThaiDate(*r.DeliveryDate)
	}
	raw := map[string]string{
		"BILLNUMBER":     r.ReceiptNO,
		"BILLDATE":       receiptDate,
		"NAME":           r.CustomerName,
		"ADDRESS2":       address,
		"DETAIL":         r.Detail,
		"DELIVERYNUMBER": deliveryNO,
		"DELIVERYDATE":   deliveryDate,
		"MONEY":          fmt.Sprintf("%.2f", r.Amount),
	}
	return raw
}

type CreateReceiptParams struct {
	TemplatePath       string
	ControlPath        string
	Filename           string
	OutputDir          string
	ReceiptDate        *string // nullable
	OverwriteReceiptNO *int    // nullable
	CustomerName       string
	CustomerID         *uint64 // nullable
	Address            *string // nullable
	Detail             string
	DeliveryNO         *string // nullable
	DeliveryDate       *string // nullable
	Amount             float64
}

func (a *App) CreateReceipt(params CreateReceiptParams) (outputPath string, err error) {
	// for logging
	logger := log.NewUnitOfLog(a.repo)
	logger.Log("CreateReceipt", params)
	// 	process data
	var receiptDate, deliveryDate *time.Time
	if receiptDate, err = utils.ParseTime(params.ReceiptDate); err != nil {
		return outputPath, logger.NewErrorAndLog(err, "ParseTime(params.ReceiptDate)")
	}
	if deliveryDate, err = utils.ParseTime(params.DeliveryDate); err != nil {
		return outputPath, logger.NewErrorAndLog(err, "ParseTime(params.DeliveryDate)")
	}

	var realReceiptNO *int
	var overwriteRowCoordinate *int
	if params.OverwriteReceiptNO != nil {
		// find original data using this NO
		originalData, err := excel.GetControlData(params.ControlPath, *params.OverwriteReceiptNO)
		if err != nil {
			return outputPath, logger.NewErrorAndLog(err, "GetControlData")
		}
		logger.Log("GetControlData", params.ControlPath, *params.OverwriteReceiptNO, originalData)
		if originalData != nil {
			realReceiptNO = &originalData.NO
			overwriteRowCoordinate = &originalData.RowCoordinate
		}
	}
	if realReceiptNO == nil {
		// get control number
		newNO, err := a.GetNextControlNumber(params.ControlPath)
		if err != nil {
			return outputPath, logger.NewErrorAndLog(err, "GetNextControlNumber")
		}
		logger.Log("GetNextControlNumber", params.ControlPath, newNO)
		realReceiptNO = &newNO
	}

	// create excel file
	receiptData := ReceiptData{
		ReceiptNO:    strconv.Itoa(*realReceiptNO),
		ReceiptDate:  receiptDate,
		CustomerName: params.CustomerName,
		Address:      params.Address,
		Detail:       params.Detail,
		DeliveryNO:   params.DeliveryNO,
		DeliveryDate: deliveryDate,
		Amount:       params.Amount,
	}
	excelFile, err := excel.CreateExcelFile(params.TemplatePath, receiptData.toExcelKeyValue())
	if err != nil {
		return outputPath, logger.NewErrorAndLog(err, "CreateExcelFile")
	}
	logger.Log("CreateReceiptFile", params.TemplatePath, receiptData)

	controlData := excel.ControlData{
		NO:           *realReceiptNO,
		CustomerName: params.CustomerName,
		Detail:       params.Detail,
		Amount:       params.Amount,
		Date:         receiptDate,
	}

	controlFile, err := excel.WriteControlFile(excel.WriteControlFileParam{
		Data:                   controlData,
		ControlFilePath:        params.ControlPath,
		OverwriteRowCoordinate: overwriteRowCoordinate,
	})
	if err != nil {
		return outputPath, logger.NewErrorAndLog(err, "WriteControlFile")
	}
	logger.Log("WriteControlFile", params.ControlPath, controlData)

	// save
	outputPath, err = excel.SaveAsExcelFile(excelFile, params.OutputDir, params.Filename)
	if err != nil {
		return outputPath, logger.NewErrorAndLog(err, "SaveAsExcelFile")
	}
	logger.Log("SaveAsExcelFile", outputPath)

	_, err = excel.SaveExcelFile(controlFile)
	if err != nil {
		return outputPath, logger.NewErrorAndLog(err, "SaveControlFile")
	}
	logger.Log("SaveControlFile", controlFile.Path)

	// insert to database if the customer does not exist
	if params.CustomerID == nil {
		newCustomer := model.Customer{Name: params.CustomerName, Address: params.Address}
		err = a.repo.CreateCustomer(&newCustomer)
		if err == nil {
			// not throwing any error!
			logger.Log("CreateCustomer", newCustomer)
		}
	}

	// open file
	err = a.CMDOpenFile(outputPath)
	if err != nil {
		return outputPath, logger.NewErrorAndLog(err, "CMDOpenFile")
	}
	return outputPath, nil
}

type ProcurementData struct {
	DeliveryNO      string
	DeliveryDate    *time.Time // nullable
	Buy             string
	Project         *string // nullable
	Amount          float64
	CustomerName    string
	Address         *string
	HeadCheckerName *string // nullable
	Checker1Name    *string // nullable
	Checker2Name    *string // nullable
	ObjectName      *string // nullable
	HeadObjectName  *string // nullable
	BossName        *string // nullable
}

func (p ProcurementData) toExcelKeyValue() map[string]string {
	var deliveryDate, project,
		address, headCheckerName, checker1Name,
		checker2Name, objectName, headObjectName, bossName string
	if p.DeliveryDate == nil {
		deliveryDate = strings.Repeat(".", 35)
	} else {
		deliveryDate = utils.GetFullThaiDate(*p.DeliveryDate)
	}
	if p.Project == nil {
		project = strings.Repeat(".", 43)
	} else {
		project = *p.Project
	}
	if p.Address != nil {
		address = *p.Address
	}
	if p.HeadCheckerName != nil {
		headCheckerName = *p.HeadCheckerName
	}
	if p.Checker1Name != nil {
		checker1Name = *p.Checker1Name
	}
	if p.Checker2Name != nil {
		checker2Name = *p.Checker2Name
	}
	if p.ObjectName != nil {
		objectName = *p.ObjectName
	}
	if p.HeadObjectName != nil {
		headObjectName = *p.HeadObjectName
	}
	if p.BossName != nil {
		bossName = *p.BossName
	}

	raw := map[string]string{
		"DELIVERYNUMBER": p.DeliveryNO,
		"DELIVERYDATE":   deliveryDate,
		"BUY":            p.Buy,
		"PROJECT":        project,
		"MONEY":          fmt.Sprintf("%.2f", p.Amount),
		"NAME":           p.CustomerName,
		"ADDRESS2":       address,
		"HEADCHECKER":    headCheckerName,
		"CHECKER1":       checker1Name,
		"CHECKER2":       checker2Name,
		"_OBJECT":        objectName,
		"OBJECT":         objectName,
		"HEADOBJECT":     headObjectName,
		"BOSS":           bossName,
	}
	return raw
}

type ProcurementOutputType string

const (
	FULL               ProcurementOutputType = "FULL"
	ONLY_DELIVERY_NOTE ProcurementOutputType = "ONLY_DELIVERY_NOTE"
	ONLY_QUOTATION     ProcurementOutputType = "ONLY_QUOTATION"
)

type CreateProcurementParams struct {
	TemplatePath          string
	ControlPath           string
	OverwriteDeliveryNO    *int
	BookOrderPath         *string // nullable
	Filename              string
	OutputDir             string
	DeliveryDate          *string // nullable
	Buy                   string
	Project               *string // nullable
	Amount                float64
	ProcurementOutputType ProcurementOutputType
	CustomerName          string
	CustomerID            *uint64 // nulllable
	Address               *string // nullable
	HeadCheckerName       *string // nullable
	Checker1Name          *string // nullable
	Checker2Name          *string // nullable
	ObjectName            *string // nullable
	HeadObjectName        *string // nullable
	BossName              *string // nullable
}

func (a *App) CreateProcurement(params CreateProcurementParams) (outputPath string, err error) {
	// for logging
	logger := log.NewUnitOfLog(a.repo)
	logger.Log("CreateProcurement", params)
	// 	process data
	var deliveryDate *time.Time
	if deliveryDate, err = utils.ParseTime(params.DeliveryDate); err != nil {
		return outputPath, logger.NewErrorAndLog(err, "ParseTime(params.DeliveryDate)")
	}
	var headCheckerName, checker1Name, checker2Name, objectName, headObjectName, bossName *string
	emptyName := strings.Repeat(" ", 35)
	headCheckerName = utils.IfNilReturnStr(params.HeadCheckerName, emptyName)
	checker1Name = utils.IfNilReturnStr(params.Checker1Name, emptyName)
	checker2Name = utils.IfNilReturnStr(params.Checker2Name, emptyName)
	objectName = utils.IfNilReturnStr(params.ObjectName, emptyName)
	headObjectName = utils.IfNilReturnStr(params.HeadObjectName, emptyName)
	bossName = utils.IfNilReturnStr(params.BossName, emptyName)

	var realDeliNO *int
	var overwriteRowCoordinate *int
	if params.OverwriteDeliveryNO != nil {
		// find original data using this NO
		originalData, err := excel.GetControlData(params.ControlPath, *params.OverwriteDeliveryNO)
		if err != nil {
			return outputPath, logger.NewErrorAndLog(err, "GetControlData")
		}
		logger.Log("GetControlData", params.ControlPath, *params.OverwriteDeliveryNO, originalData)
		if originalData != nil {
			realDeliNO = &originalData.NO
			overwriteRowCoordinate = &originalData.RowCoordinate
		}
	}
	if realDeliNO == nil {
		// get control number
		newNO, err := a.GetNextControlNumber(params.ControlPath)
		if err != nil {
			return outputPath, logger.NewErrorAndLog(err, "GetNextControlNumber")
		}
		logger.Log("GetNextControlNumber", params.ControlPath, newNO)
		realDeliNO = &newNO
	}

	// create excel file
	procurementData := ProcurementData{
		DeliveryNO:      strconv.Itoa(*realDeliNO),
		DeliveryDate:    deliveryDate,
		Buy:             params.Buy,
		Project:         params.Project,
		Amount:          params.Amount,
		CustomerName:    params.CustomerName,
		Address:         params.Address,
		HeadCheckerName: headCheckerName,
		Checker1Name:    checker1Name,
		Checker2Name:    checker2Name,
		ObjectName:      objectName,
		HeadObjectName:  headObjectName,
		BossName:        bossName,
	}
	excelFile, err := excel.CreateExcelFile(params.TemplatePath, procurementData.toExcelKeyValue())
	if err != nil {
		return outputPath, logger.NewErrorAndLog(err, "CreateExcelFile")
	}
	logger.Log("CreateExcelFile", params.TemplatePath, procurementData)

	// book order
	if params.BookOrderPath != nil {
		bookOrder, err := excel.GetBookOrderDataFromFile(*params.BookOrderPath)
		if err != nil {
			return outputPath, logger.NewErrorAndLog(err, "GetBookOrderDataFromFile")
		}
		logger.Log("GetBookOrderDataFromFile", *params.BookOrderPath, bookOrder)
		excelFile, err = excel.WriteBookOrder(excelFile, "data", bookOrder)
		if err != nil {
			return outputPath, logger.NewErrorAndLog(err, "WriteBookOrder")
		}
		logger.Log("WriteBookOrder", bookOrder)
	}

	switch params.ProcurementOutputType {
	case ONLY_DELIVERY_NOTE:
		excelFile, err = excel.ShowOnlySheetNames(excelFile, "ใบส่งของ")
		if err != nil {
			return outputPath, logger.NewErrorAndLog(err, "ShowOnlySheetNames(ONLY_DELIVERY_NOTE)")
		}
	case ONLY_QUOTATION:
		excelFile, err = excel.ShowOnlySheetNames(excelFile, "ใบเสนอราคา", "แนบใบเสนอราคา")
		if err != nil {
			return outputPath, logger.NewErrorAndLog(err, "ShowOnlySheetNames(ONLY_QUOTATION)")
		}
	}

	controlData := excel.ControlData{
		NO:           *realDeliNO,
		CustomerName: params.CustomerName,
		Detail:       params.Buy,
		Amount:       params.Amount,
		Date:         deliveryDate,
	}
	controlFile, err := excel.WriteControlFile(excel.WriteControlFileParam{
		Data:                   controlData,
		ControlFilePath:        params.ControlPath,
		OverwriteRowCoordinate: overwriteRowCoordinate,
	})
	if err != nil {
		return outputPath, logger.NewErrorAndLog(err, "WriteControlFile")
	}
	logger.Log("WriteControlFile", params.ControlPath, controlData)

	// save
	outputPath, err = excel.SaveAsExcelFile(excelFile, params.OutputDir, params.Filename)
	if err != nil {
		return outputPath, logger.NewErrorAndLog(err, "SaveAsExcelFile")
	}
	logger.Log("SaveAsExcelFile", outputPath)

	_, err = excel.SaveExcelFile(controlFile)
	if err != nil {
		return outputPath, logger.NewErrorAndLog(err, "SaveControlFile")
	}
	logger.Log("SaveControlFile", controlFile.Path)

	// upsert
	customer := model.Customer{
		Name:            params.CustomerName,
		Address:         params.Address,
		HeadCheckerName: params.HeadCheckerName,
		Checker1Name:    params.Checker1Name,
		Checker2Name:    params.Checker2Name,
		ObjectName:      params.ObjectName,
		HeadObjectName:  params.HeadObjectName,
		BossName:        params.BossName,
	}
	if params.CustomerID == nil {
		err = a.repo.CreateCustomer(&customer)
		if err == nil {
			// not throwing any error!
			logger.Log("CreateCustomer", customer)
		}
	} else {
		customer.ID = *params.CustomerID
		_, err = a.repo.UpdateCustomerByID(&customer)
		if err == nil {
			// not throwing any error!
			logger.Log("UpdateCustomerByID", customer)
		}
	}

	// open file
	err = a.CMDOpenFile(outputPath)
	if err != nil {
		return outputPath, logger.NewErrorAndLog(err, "CMDOpenFile")
	}
	return outputPath, nil
}

func (a *App) GetBookOrderFromDataSourceFile(filePath string) (excel.BookOrder, error) {
	logger := log.NewSingleLog(a.repo)
	data, err := excel.GetBookOrderDataFromFile(filePath)
	if err != nil {
		return nil, logger.NewErrorAndLog(err, "GetBookOrderFromDataSourceFile")
	}
	logger.Log("GetBookOrderFromDataSourceFile", filePath)
	return data, nil
}

func (a *App) AutoMoveBookOrder(procurementFilepath string, bookOrderFilePath string) error {
	logger := log.NewUnitOfLog(a.repo)
	logger.Log("AutoMoveBookOrder", procurementFilepath, bookOrderFilePath)
	f, err := excel.AutoMoveBookOrder(procurementFilepath, bookOrderFilePath)
	if err != nil {
		return logger.NewErrorAndLog(err, "AutoMoveBookOrder")
	}
	logger.Log("AutoMoveBookOrder", procurementFilepath, bookOrderFilePath)

	_, err = excel.SaveExcelFile(f)
	if err != nil {
		return logger.NewErrorAndLog(err, "SaveExcelFile")
	}

	err = a.CMDOpenFile(procurementFilepath)
	if err != nil {
		return logger.NewErrorAndLog(err, "CMDOpenFile")
	}
	return nil
}
