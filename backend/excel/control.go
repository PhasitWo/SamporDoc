package excel

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ControlData struct {
	NO           int
	CustomerName string
	Detail       string
	Amount       float64
	Date         *time.Time
}

type ControlDataWithRowCoordinate struct {
	ControlData
	RowCoordinate int
}

type WriteControlFileParam struct {
	ControlFilePath        string
	Data                   ControlData
	OverwriteRowCoordinate *int
}

func WriteControlFile(param WriteControlFileParam) (*excelize.File, error) {
	data := param.Data
	controlFilePath := param.ControlFilePath
	overwriteRowCoordinate := param.OverwriteRowCoordinate
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

	var toInsertRowCoordinate int
	const startRow = 3
	activeSheetName := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.Rows(activeSheetName)
	if overwriteRowCoordinate != nil {
		toInsertRowCoordinate = *overwriteRowCoordinate
	} else {
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
		if rowCoodinate < startRow {
			toInsertRowCoordinate = startRow
		} else if hasGap {
			toInsertRowCoordinate = rowCoodinate
		} else {
			toInsertRowCoordinate = rowCoodinate + 1
		}
	}

	f.SetCellValue(activeSheetName, getCellName(1, toInsertRowCoordinate), data.NO)
	f.SetCellValue(activeSheetName, getCellName(2, toInsertRowCoordinate), data.CustomerName)
	f.SetCellValue(activeSheetName, getCellName(3, toInsertRowCoordinate), data.Detail)
	f.SetCellValue(activeSheetName, getCellName(4, toInsertRowCoordinate), data.Amount)
	f.SetCellValue(activeSheetName, getCellName(5, toInsertRowCoordinate), date.Format("02-01-2006"))

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

	// collect all number
	rowCoodinate := 0
	numSlice := []int{}
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
		NO, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			return -1, newError(err)
		}
		if NO > 0 {
			numSlice = append(numSlice, NO)
		}
	}
	if len(numSlice) == 0 {
		return 1, nil
	}
	// sort asc
	slices.SortFunc(numSlice, func(a int, b int) int {
		return a - b
	})
	candidateNo := -1
	// determine next number
	for i := range math.MaxInt {
		expected := i + 1
		if i >= len(numSlice) {
			candidateNo = expected
			break
		}
		if numSlice[i] != expected {
			candidateNo = expected
			break
		}
	}
	return candidateNo, nil
}

func GetControlData(controlFilePath string, NO int) (*ControlDataWithRowCoordinate, error) {
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
	rowCoodinate := 0
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
			continue
		}
		curNO, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			continue
		}
		if curNO == NO {
			// found the control data
			customerName, _ := f.GetCellValue(activeSheetName, getCellName(2, rowCoodinate))
			detail, _ := f.GetCellValue(activeSheetName, getCellName(3, rowCoodinate))
			rawAmount, _ := f.GetCellValue(activeSheetName, getCellName(4, rowCoodinate), excelize.Options{RawCellValue: true})
			amount, err := strconv.ParseFloat(strings.TrimSpace(rawAmount), 64)
			if err != nil {
				amount = 0
			}
			return &ControlDataWithRowCoordinate{
				ControlData: ControlData{
					NO:           curNO,
					CustomerName: customerName,
					Detail:       detail,
					Amount:       amount,
				},
				RowCoordinate: rowCoodinate,
			}, nil
		}
	}
	return nil, nil
}
