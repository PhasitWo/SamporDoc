package excel

import (
	"SamporDoc/backend/utils"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type BookOrder = []PublisherItem

type PublisherItem struct {
	Name        string
	Items       []BookItem
	TotalAmount float64
}

type BookItem struct {
	Name          string
	Quantity      int
	Price         float64
	PublisherName string
}

func GetBookOrderDataFromFile(filepath string) (BookOrder, error) {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return nil, newError(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	PUBLISHER_CELL_ADDRESS := "F4"
	data := []PublisherItem{}
	dataMap := make(map[string]*PublisherItem)
	sheetList := f.GetSheetList()
	for _, activeSheetName := range sheetList {
		rows, err := f.GetRows(activeSheetName)
		if err != nil {
			return nil, newError(err)
		}
		publisherName, _ := f.GetCellValue(activeSheetName, PUBLISHER_CELL_ADDRESS)
		if strings.TrimSpace(publisherName) == "" {
			publisherName = activeSheetName
		}
		targetPublisher := dataMap[publisherName]
		if targetPublisher == nil {
			targetPublisher = &PublisherItem{Name: publisherName}
			dataMap[publisherName] = targetPublisher
		}
		for ridx, row := range rows {
			for cidx, cell := range row {
				adjCell1, _ := f.GetCellValue(activeSheetName, getCellName(cidx+2, ridx+1))
				adjCell2, _ := f.GetCellValue(activeSheetName, getCellName(cidx+3, ridx+1))
				cleanCell := strings.TrimSpace(cell)
				cleanAdjCell1 := strings.TrimSpace(adjCell1)
				cleanAdjCell2 := strings.TrimSpace(adjCell2)
				if !utils.IsFloat(cleanCell) && utils.IsFloat(cleanAdjCell1) && utils.IsFloat(cleanAdjCell2) {
					price, _ := strconv.ParseFloat(cleanAdjCell1, 64)
					quantity, _ := strconv.Atoi(cleanAdjCell2)
					targetPublisher.Items = append(targetPublisher.Items, BookItem{
						Name:          cleanCell,
						Quantity:      quantity,
						Price:         price,
						PublisherName: publisherName,
					})
					// fmt.Printf("Found item: %s, price: %f, quantity: %d sum: %f\n", cleanCell, price, quantity, price*float64(quantity))
					targetPublisher.TotalAmount += price * float64(quantity)
				}

			}
		}
	}
	for _, publisher := range dataMap {
		if len(publisher.Items) > 0 {
			data = append(data, *publisher)
		}
	}
	// sort data
	slices.SortFunc(data, func(a, b PublisherItem) int {
		return int(b.TotalAmount) - int(a.TotalAmount)
	})
	return data, nil
}

func WriteBookOrder(f *excelize.File, sheetName string, bookOrder BookOrder) (*excelize.File, error) {
	idx, err := f.GetSheetIndex(sheetName)
	if idx == -1 || err != nil {
		return nil, newError(fmt.Errorf("sheet %s not found in template", sheetName))
	}

	allItems := []BookItem{}
	for _, publisher := range bookOrder {
		allItems = append(allItems, publisher.Items...)
	}
	SortBookItemArray(&allItems)

	START_ROW := 2
	rowIdx := START_ROW
	for _, item := range allItems {
		f.SetCellValue(sheetName, getCellName(2, rowIdx), item.Name)
		f.SetCellValue(sheetName, getCellName(3, rowIdx), item.Quantity)
		f.SetCellValue(sheetName, getCellName(4, rowIdx), "เล่ม")
		f.SetCellValue(sheetName, getCellName(5, rowIdx), item.Price)
		f.SetCellValue(sheetName, getCellName(7, rowIdx), item.PublisherName)
		rowIdx++
	}
	return f, nil
}

type SortRule struct {
	Keyword string
	Level   int
}

var sortRules = []SortRule{
	{"นิทาน", 0}, {"อนุบาล", 1},
	{"อ.1", 2}, {"อ.2", 3}, {"อ.3", 4},
	{"ป.1", 5}, {"ป.2", 6}, {"ป.3", 7}, {"ป.4", 8}, {"ป.5", 9}, {"ป.6", 10},
	{"ม.1", 11}, {"ม.2", 12}, {"ม.3", 13}, {"ม.4-6", 17}, {"ม.4", 14}, {"ม.5", 15}, {"ม.6", 16},
}

func getSortLevel(name string) int {
	cleanName := strings.ReplaceAll(name, " ", "")
	lastLevel := 9999
	lastIndex := -1
	for _, rule := range sortRules {
		idx := strings.LastIndex(cleanName, rule.Keyword)
		if idx != -1 && idx > lastIndex {
			lastIndex = idx
			lastLevel = rule.Level
		}
	}
	return lastLevel
}

func SortBookItemArray(items *[]BookItem) {
	slices.SortFunc(*items, func(a, b BookItem) int {
		return getSortLevel(a.Name) - getSortLevel(b.Name)
	})
}

func AutoMoveBookOrder(procurementFilepath string, bookOrderFilePath string) (*excelize.File, error) {
	f, err := excelize.OpenFile(procurementFilepath)
	if err != nil {
		return nil, newError(err)
	}
	bookOrder, err := GetBookOrderDataFromFile(bookOrderFilePath)
	if err != nil {
		return nil, newError(err)
	}
	f, err = WriteBookOrder(f, "data", bookOrder)
	if err != nil {
		return nil, newError(err)
	}
	return f, nil
}
