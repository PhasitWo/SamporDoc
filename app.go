package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"SamporDoc/backend/config"
	"SamporDoc/backend/excel"
	"SamporDoc/backend/model"
	"SamporDoc/backend/repository"
	"SamporDoc/backend/seed"
	"SamporDoc/backend/utils"

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
	db.AutoMigrate(&model.School{}, &model.Shop{}, &model.Log{})

	// create repo instance
	a.repo = repository.NewRepo(ctx, db)

	// seeding
	schools, err := a.repo.GetAllSchools()
	if err != nil {
		fmt.Println("Error seeding database -> repo.getAllSchools")
		panic(err)
	}
	if len(schools) == 0 {
		seed.SeedSchools(a.repo)
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

func (a *App) UpdateShopBySlug(shop model.Shop) error {
	return a.repo.UpdateShopBySlug(&shop)
}

func (a *App) GetAllSchools() ([]model.School, error) {
	return a.repo.GetAllSchools()
}

func (a *App) GetNextControlNumber(controlFilePath string) (int, error) {
	return excel.GetNextControlNumber(controlFilePath)
}

type ReceiptData struct {
	ReceiptNO    int
	ReceiptDate  *time.Time
	CustomerName string
	Address      *string
	Detail       *string
	DeliveryNO   *int
	DeliveryDate *time.Time
	Amount       float64
}

func (r ReceiptData) toExcelKeyValue() map[string]string {
	raw := map[string]string{
		"BILLNUMBER":     strconv.Itoa(r.ReceiptNO),
		"BILLDATE":       utils.IfNilReturnStr(r.ReceiptDate, strings.Repeat(".", 40), utils.GetFullThaiDate(*r.ReceiptDate)),
		"NAME":           r.CustomerName,
		"ADDRESS2":       utils.IfNilReturnStr(r.Address, "", *r.Address),
		"DETAIL":         utils.IfNilReturnStr(r.Detail, "", *r.Detail),
		"DELIVERYNUMBER": utils.IfNilReturnStr(r.DeliveryNO, "", utils.GetShortThaiDate(*r.DeliveryDate)),
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
