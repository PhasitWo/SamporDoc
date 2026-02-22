package excel

import (
	"SamporDoc/backend/utils"
	"fmt"
	"maps"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func newError(err error) error {
	return fmt.Errorf("[EXCEL]: %w", err)
}

func SaveExcelFile(f *excelize.File) error {
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	err := f.Save()
	if err != nil {
		return newError(err)
	}

	return nil
}

func SaveAsExcelFile(f *excelize.File, outputDir string, filename string) (string, error) {
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// long path workaround for windows
	outputPath := filepath.Join(outputDir, filename+".xlsx")
	if runtime.GOOS == "windows" && len(outputPath) > 180 {
		drive, err := utils.GetAvailableDriveLetter()
		if err != nil {
			return "", newError(err)
		}

		if err := utils.MapDrive(drive, outputDir); err != nil {
			return "", newError(err)
		}
		defer utils.UnmapDrive(drive)

		shortPath := fmt.Sprintf("%s\\%s", drive, filename+".xlsx")

		if err := f.SaveAs(shortPath); err != nil {
			return "", newError(err)
		}
	} else {
		// for non-windows or short path
		err := f.SaveAs(outputPath)
		if err != nil {
			return "", newError(err)
		}
	}
	return outputPath, nil
}

func CreateExcelFile(templatePath string, data map[string]string) (*excelize.File, error) {
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return nil, newError(err)
	}

	sheetList := f.GetSheetList()
	sortedKeys := slices.SortedFunc(maps.Keys(data), func(a, b string) int {
		if len(b) != len(a) {
			return len(b) - len(a)
		}
		return strings.Compare(a, b) // Alphabetical if lengths are equal
	})
	for _, activeSheetName := range sheetList {
		rows, err := f.GetRows(activeSheetName)
		if err != nil {
			return nil, newError(err)
		}
		for ridx, row := range rows {
			for cidx, cell := range row {
				coordinate, _ := excelize.CoordinatesToCellName(cidx+1, ridx+1)
				newValue := cell
				replaced := false
				for _, key := range sortedKeys {
					replacement := data[key]
					if strings.Contains(newValue, key) {
						replaced = true
						newValue = strings.ReplaceAll(newValue, key, replacement)
					}
				}
				if replaced {
					if numeric, err := strconv.ParseFloat(newValue, 64); err == nil {
						f.SetCellValue(activeSheetName, coordinate, numeric)
					} else {
						f.SetCellValue(activeSheetName, coordinate, newValue)
					}
				}

			}
		}
	}
	return f, nil
}

type ControlData struct {
	NO           int
	CustomerName string
	Detail       string
	Amount       float64
	Date         *time.Time
}

func getCellName(col int, row int) string {
	cellName, _ := excelize.CoordinatesToCellName(col, row)
	return cellName
}

func WriteControlFile(controlFilePath string, data ControlData) (*excelize.File, error) {
	// process date
	date := time.Now()
	if data.Date != nil {
		date = *data.Date
	}
	// excel
	f, err := excelize.OpenFile(controlFilePath)
	if err != nil {
		return nil, newError(err)
	}

	const startRow = 3
	activeSheetName := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.Rows(activeSheetName)
	// find empty cell to insert
	rowCoodinate := 0
	hasGap := false
	for rows.Next() {
		rowCoodinate++
		if rowCoodinate < startRow {
			continue
		}
		row, err := rows.Columns()
		if err != nil {
			return nil, newError(err)
		}
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			hasGap = true
			break
		}
	}
	var toInsertRowCoordinate int
	if rowCoodinate < startRow {
		toInsertRowCoordinate = startRow
	} else if hasGap {
		toInsertRowCoordinate = rowCoodinate
	} else {
		toInsertRowCoordinate = rowCoodinate + 1
	}

	f.SetCellValue(activeSheetName, getCellName(1, toInsertRowCoordinate), data.NO)
	f.SetCellValue(activeSheetName, getCellName(2, toInsertRowCoordinate), data.CustomerName)
	f.SetCellValue(activeSheetName, getCellName(3, toInsertRowCoordinate), data.Detail)
	f.SetCellValue(activeSheetName, getCellName(4, toInsertRowCoordinate), data.Amount)
	f.SetCellValue(activeSheetName, getCellName(5, toInsertRowCoordinate), date.Format("02-01-2006"))

	return f, nil
}

func ShowOnlySheetNames(f *excelize.File, sheetNames ...string) (*excelize.File, error) {
	for _, name := range f.GetSheetList() {
		if slices.Contains(sheetNames, name) {
			idx, _ := f.GetSheetIndex(name)
			f.SetActiveSheet(idx)
			break
		}
	}
	for _, name := range f.GetSheetList() {
		if !slices.Contains(sheetNames, name) {
			err := f.SetSheetVisible(name, false)
			if err != nil {
				return nil, err
			}
		}
	}
	return f, nil
}

// find maximun number regardless of gap
func GetNextControlNumber(controlFilePath string) (int, error) {
	f, err := excelize.OpenFile(controlFilePath)
	if err != nil {
		return -1, newError(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	const startRow = 3
	activeSheetName := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.Rows(activeSheetName)
	// find last row
	rowCoodinate := 0
	maxNO := 0
	for rows.Next() {
		rowCoodinate++
		if rowCoodinate < startRow {
			continue
		}
		row, err := rows.Columns()
		if err != nil {
			return -1, newError(err)
		}
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}
		maxNO, err = strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return -1, newError(err)
		}
	}
	return maxNO + 1, nil
}

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

	// PUBLISHER_CELL_ADDRESS := "F4"
	data := []PublisherItem{}
	dataMap := make(map[string]*PublisherItem)
	sheetList := f.GetSheetList()
	for _, activeSheetName := range sheetList {
		rows, err := f.GetRows(activeSheetName)
		if err != nil {
			return nil, newError(err)
		}
		// publisherName, _ := f.GetCellValue(activeSheetName, PUBLISHER_CELL_ADDRESS)
		publisherName := activeSheetName
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

type SortRule struct {
	Keyword string
	Level   int
}

var sortRules = []SortRule{
	{"นิทาน", 0}, {"อนุบาล", 1},
	{"อ.1", 2}, {"อ.2", 3}, {"อ.3", 4},
	{"ป.1", 5}, {"ป.2", 6}, {"ป.3", 7}, {"ป.4", 8}, {"ป.5", 9}, {"ป.6", 10},
	{"ม.1", 11}, {"ม.2", 12}, {"ม.3", 13}, {"ม.4", 14}, {"ม.5", 15}, {"ม.6", 16},
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

func getSortLevel(name string) int {
	cleanName := strings.ReplaceAll(name, " ", "")
	for _, rule := range sortRules {
		if strings.Contains(cleanName, rule.Keyword) {
			return rule.Level
		}
	}
	return 9999
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
