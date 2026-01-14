package excel

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func newError(err error) error {
	return fmt.Errorf("[EXCEL]: %w", err)
}

func CreateExcelFile(templatePath string, data map[string]string) (*excelize.File, error) {
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return nil, newError(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
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
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
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
