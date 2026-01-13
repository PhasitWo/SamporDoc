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
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
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
	Detail       *string
	DeliveryNO   *string
	DeliveryDate *time.Time
	Amount       float64
}

func (r ReceiptData) toExcelKeyValue() map[string]string {
	var receiptDate, address, detail, deliveryNO, deliveryDate string
	if r.ReceiptDate == nil {
		receiptDate = strings.Repeat(".", 40)
	} else {
		receiptDate = utils.GetFullThaiDate(*r.ReceiptDate)
	}
	if r.Address != nil {
		address = *r.Address
	}
	if r.Detail != nil {
		detail = *r.Detail
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
		"DETAIL":         detail,
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
	Address      *string // nullable
	Detail       *string // nullable
	DeliveryNO   *string // nullable
	DeliveryDate *string // nullable
	Amount       float64
}

func (a *App) CreateReceipt(params CreateReceiptParams) (err error) {
	// for logging
	logger := log.NewUnitOfLog(a.repo)
	// 	process data
	var receiptDate, deliveryDate *time.Time
	if receiptDate, err = utils.ParseTime(params.ReceiptDate); err != nil {
		return logger.NewErrorAndLog(err, "ParseTime(params.ReceiptDate)")
	}
	if deliveryDate, err = utils.ParseTime(params.DeliveryDate); err != nil {
		return logger.NewErrorAndLog(err, "ParseTime(params.DeliveryDate)")
	}
	// insert to database if the customerName does not exist
	if params.Address != nil {
		_, err = a.repo.GetCustomerByName(params.CustomerName)
		if errors.Is(errors.Unwrap(err), gorm.ErrRecordNotFound) {
			// not throwing any error!
			newCustomer := model.Customer{Name: params.CustomerName, Address: params.Address}
			err = a.repo.CreateCustomer(&newCustomer)
			if err == nil {
				logger.Log("CreateCustomer", newCustomer)
			}
		}
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
		return logger.NewErrorAndLog(err, "CreateExcelFile")
	}
	logger.Log("CreateExcelFile", params.TemplatePath, receiptData)

	controlData := excel.ControlData{
		NO:           params.ReceiptNO,
		CustomerName: params.CustomerName,
		Amount:       params.Amount,
		Date:         receiptDate,
	}
	controlFile, err := excel.WriteControlFile(params.ControlPath, controlData)
	if err != nil {
		return logger.NewErrorAndLog(err, "WriteControlFile")
	}
	logger.Log("WriteControlFile", params.ControlPath, controlData)

	// save
	outputFilePath := filepath.Join(params.OutputDir, params.Filename+".xlsx")
	err = excelFile.SaveAs(outputFilePath)
	if err != nil {
		return logger.NewErrorAndLog(err, "SaveOuputReceiptFile")
	}
	logger.Log("SaveOuputReceiptFile", outputFilePath)

	err = controlFile.Save()
	if err != nil {
		return logger.NewErrorAndLog(err, "SaveControlFile")
	}
	logger.Log("SaveControlFile", controlFile.Path)

	// open file
	err = a.CMDOpenFile(outputFilePath)
	if err != nil {
		return logger.NewErrorAndLog(err, "SaveControlFile")
	}
	return nil
}
