package excel

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func newError(err error) error {
	return fmt.Errorf("[EXCEL]: %w", err)
}

type a = map[string]string

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
	for _, activeSheetName := range sheetList {
		rows, err := f.GetRows(activeSheetName)
		if err != nil {
			return nil, newError(err)
		}
		for range 3 {
			for ridx, row := range rows {
				for cidx, cell := range row {
					coordinate, _ := excelize.CoordinatesToCellName(cidx+1, ridx+1)
					for key, replacement := range data {
						if strings.Contains(cell, key) {
							newValue := strings.ReplaceAll(cell, key, replacement)
							f.SetCellValue(activeSheetName, coordinate, newValue)
						}
					}
				}
			}
		}
	}
	return f, nil
}

type ControlData struct {
	NO           int
	customerName string
	amount       int
	date         time.Time
}

func getCellName(col int, row int) string {
	cellName, _ := excelize.CoordinatesToCellName(col, row)
	return cellName
}

func WriteControlFile(controlFilePath string, data ControlData) (*excelize.File, error) {
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
	// find last row
	rowCoodinate := 0
	prev := ""
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
			break
		}
		prev = strings.TrimSpace(row[0])
	}

	if rowCoodinate < startRow || prev == "" {
		f.SetCellValue(activeSheetName, getCellName(1, startRow), 1)
		f.SetCellValue(activeSheetName, getCellName(2, startRow), data.customerName)
		f.SetCellValue(activeSheetName, getCellName(3, startRow), data.amount)
		f.SetCellValue(activeSheetName, getCellName(4, startRow), data.date.Format("02-01-2006"))
	} else {
		f.SetCellValue(activeSheetName, getCellName(1, rowCoodinate), data.NO)
		f.SetCellValue(activeSheetName, getCellName(2, rowCoodinate), data.customerName)
		f.SetCellValue(activeSheetName, getCellName(3, rowCoodinate), data.amount)
		f.SetCellValue(activeSheetName, getCellName(4, rowCoodinate), data.date.Format("02-01-2006"))
	}
	return f, nil
}

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
	prev := ""
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
			break
		}
		prev = strings.TrimSpace(row[0])
	}
	if prev == "" {
		return 1, nil
	}
	latestValue, err := strconv.Atoi(prev)
	if err != nil {
		return -1, newError(err)
	}
	return latestValue + 1, nil
}
