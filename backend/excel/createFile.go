package excel

import (
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

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
