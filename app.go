package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"SamporDoc/backend/config"
	"SamporDoc/backend/excel"
	"SamporDoc/backend/model"
	"SamporDoc/backend/repository"
	"SamporDoc/backend/seed"
	"SamporDoc/backend/utils"

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
	raw := map[string]string{
		"BILLNUMBER":     r.ReceiptNO,
		"BILLDATE":       utils.IfNilReturnStr(r.ReceiptDate, strings.Repeat(".", 40), utils.GetFullThaiDate(*r.ReceiptDate)),
		"NAME":           r.CustomerName,
		"ADDRESS2":       utils.IfNilReturnStr(r.Address, "", *r.Address),
		"DETAIL":         utils.IfNilReturnStr(r.Detail, "", *r.Detail),
		"DELIVERYNUMBER": utils.IfNilReturnStr(r.DeliveryNO, "", *r.DeliveryNO),
		"DELIVERYDATE":   utils.IfNilReturnStr(r.DeliveryDate, "", utils.GetShortThaiDate(*r.DeliveryDate)),
		"MONEY":          fmt.Sprintf("%.2f", r.Amount),
	}
	return raw
}

type CreateReceiptParams struct {
	TemplatePath   string
	OutputFilePath string
	Data           ReceiptData
	ControlPath    string
	ControlData    excel.ControlData
}

func (a *App) CreateReceipt(params CreateReceiptParams) error {
	excelFile, err := excel.CreateExcelFile(params.TemplatePath, params.Data.toExcelKeyValue())
	if err != nil {
		return err
	}
	controlFile, err := excel.WriteControlFile(params.ControlPath, params.ControlData)
	if err != nil {
		return err
	}
	// save
	err = excelFile.SaveAs(params.OutputFilePath)
	if err != nil {
		return err
	}
	err = controlFile.Save()
	if err != nil {
		return err
	}
	return nil
}
