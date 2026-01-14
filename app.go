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
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"gorm.io/driver/sqlite"
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
	}
}

func (a *App) CMDOpenFile(filePath string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", filePath)
	} else {
		cmd = exec.Command("cmd", "/C", "start", filePath)
	}
	return cmd.Run()
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
	TemplatePath string
	ControlPath  string
	Filename     string
	OutputDir    string
	ReceiptNO    string
	ReceiptDate  *string // nullable
	CustomerName string
	CustomerID   *uint64 // nullable
	Address      *string // nullable
	Detail       string
	DeliveryNO   *string // nullable
	DeliveryDate *string // nullable
	Amount       float64
}

func (a *App) CreateReceipt(params CreateReceiptParams) (err error) {
	// for logging
	logger := log.NewUnitOfLog(a.repo)
	// 	process data
	var receiptDate, deliveryDate *time.Time
	var convertedReceiptNO int
	if receiptDate, err = utils.ParseTime(params.ReceiptDate); err != nil {
		return logger.NewErrorAndLog(err, "ParseTime(params.ReceiptDate)")
	}
	if deliveryDate, err = utils.ParseTime(params.DeliveryDate); err != nil {
		return logger.NewErrorAndLog(err, "ParseTime(params.DeliveryDate)")
	}
	if convertedReceiptNO, err = strconv.Atoi(params.ReceiptNO); err != nil {
		return logger.NewErrorAndLog(err, "strconv.Atoi(params.ReceiptNO)")
	}

	// create excel file
	receiptData := ReceiptData{
		ReceiptNO:    params.ReceiptNO,
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
		return logger.NewErrorAndLog(err, "CreateReceiptFile")
	}
	logger.Log("CreateReceiptFile", params.TemplatePath, receiptData)

	controlData := excel.ControlData{
		NO:           convertedReceiptNO,
		CustomerName: params.CustomerName,
		Detail:       params.Detail,
		Amount:       params.Amount,
		Date:         receiptDate,
	}
	controlFile, err := excel.WriteControlFile(params.ControlPath, controlData)
	if err != nil {
		return logger.NewErrorAndLog(err, "WriteReceiptControlFile")
	}
	logger.Log("WriteReceiptControlFile", params.ControlPath, controlData)

	// save
	outputFilePath := filepath.Join(params.OutputDir, params.Filename+".xlsx")
	err = excelFile.SaveAs(outputFilePath)
	if err != nil {
		return logger.NewErrorAndLog(err, "SaveOuputReceiptFile")
	}
	logger.Log("SaveOuputReceiptFile", outputFilePath)

	err = controlFile.Save()
	if err != nil {
		return logger.NewErrorAndLog(err, "SaveReceiptControlFile")
	}
	logger.Log("SaveReceiptControlFile", controlFile.Path)

	// insert to database if the customer does not exist
	if params.CustomerID != nil {
		newCustomer := model.Customer{Name: params.CustomerName, Address: params.Address}
		err = a.repo.CreateCustomer(&newCustomer)
		if err == nil {
			// not throwing any error!
			logger.Log("CreateCustomer", newCustomer)
		}
	}

	// open file
	err = a.CMDOpenFile(outputFilePath)
	if err != nil {
		return logger.NewErrorAndLog(err, "CMDOpenReceiptFile")
	}
	return nil
}

type ProcurementData struct {
	DeliveryNO      string
	DeliveryDate    *time.Time // nullable
	Buy             string
	Project         *string // nullable
	Amount          float64
	Quantity        *int // nullable
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
	var deliveryDate, project, quantity,
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
	if p.Quantity == nil {
		quantity = "**QUANTITY**"
	} else {
		quantity = strconv.Itoa(*p.Quantity)
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
		"QUANTITY":       quantity,
		"NAME":           p.CustomerName,
		"ADDRESS2":       address,
		"HEADCHECKER":    headCheckerName,
		"CHECKER1":       checker1Name,
		"CHECKER2":       checker2Name,
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
	Filename              string
	OutputDir             string
	DeliveryNO            string
	DeliveryDate          *string // nullable
	Buy                   string
	Project               *string // nullable
	Amount                float64
	ProcurementOutputType ProcurementOutputType
	Quantity              *int // nulllable
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

func (a *App) CreateProcurement(params CreateProcurementParams) (err error) {
	// for logging
	logger := log.NewUnitOfLog(a.repo)
	// 	process data
	var deliveryDate *time.Time
	var convertedDeliveryNO int
	if deliveryDate, err = utils.ParseTime(params.DeliveryDate); err != nil {
		return logger.NewErrorAndLog(err, "ParseTime(params.DeliveryDate)")
	}
	if convertedDeliveryNO, err = strconv.Atoi(params.DeliveryNO); err != nil {
		return logger.NewErrorAndLog(err, "strconv.Atoi(params.DeliveryNO)")
	}
	// create excel file
	procurementData := ProcurementData{
		DeliveryNO:      params.DeliveryNO,
		DeliveryDate:    deliveryDate,
		Buy:             params.Buy,
		Project:         params.Project,
		Amount:          params.Amount,
		Quantity:        params.Quantity,
		CustomerName:    params.CustomerName,
		Address:         params.Address,
		HeadCheckerName: params.HeadCheckerName,
		Checker1Name:    params.Checker1Name,
		Checker2Name:    params.Checker2Name,
		ObjectName:      params.ObjectName,
		HeadObjectName:  params.HeadObjectName,
		BossName:        params.BossName,
	}
	excelFile, err := excel.CreateExcelFile(params.TemplatePath, procurementData.toExcelKeyValue())
	if err != nil {
		return logger.NewErrorAndLog(err, "CreateProcurementFile")
	}
	logger.Log("CreateProcurementFile", params.TemplatePath, procurementData)

	switch params.ProcurementOutputType {
	case ONLY_DELIVERY_NOTE:
		excelFile, err = excel.ShowOnlySheetNames(excelFile, "ใบส่งของ")
		if err != nil {
			return logger.NewErrorAndLog(err, "ShowOnlySheetNames(ONLY_DELIVERY_NOTE)")
		}
	case ONLY_QUOTATION:
		excelFile, err = excel.ShowOnlySheetNames(excelFile, "ใบเสนอราคา", "แนบใบเสนอราคา")
		if err != nil {
			return logger.NewErrorAndLog(err, "ShowOnlySheetNames(ONLY_QUOTATION)")
		}
	}

	controlData := excel.ControlData{
		NO:           convertedDeliveryNO,
		CustomerName: params.CustomerName,
		Detail:       params.Buy,
		Amount:       params.Amount,
		Date:         deliveryDate,
	}
	controlFile, err := excel.WriteControlFile(params.ControlPath, controlData)
	if err != nil {
		return logger.NewErrorAndLog(err, "WriteProcurementControlFile")
	}
	logger.Log("WriteProcurementControlFile", params.ControlPath, controlData)

	// save
	outputFilePath := filepath.Join(params.OutputDir, params.Filename+".xlsx")
	err = excelFile.SaveAs(outputFilePath)
	if err != nil {
		return logger.NewErrorAndLog(err, "SaveOuputProcurementFile")
	}
	logger.Log("SaveOuputProcurementFile", outputFilePath)

	err = controlFile.Save()
	if err != nil {
		return logger.NewErrorAndLog(err, "SaveProcurementControlFile")
	}
	logger.Log("SaveProcurementControlFile", controlFile.Path)

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
			logger.Log("CreateCustomer", customer)
		}
	}

	// open file
	err = a.CMDOpenFile(outputFilePath)
	if err != nil {
		return logger.NewErrorAndLog(err, "CMDOpenProcurementFile")
	}
	return nil
}
