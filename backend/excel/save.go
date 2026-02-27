package excel

import (
	"SamporDoc/backend/utils"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/xuri/excelize/v2"
)

func SaveAsExcelFile(f *excelize.File, outputDir string, filename string) (string, error) {
	// force reevaluate all formulas before saving
	t := true
	f.SetCalcProps(&excelize.CalcPropsOptions{FullCalcOnLoad: &t, ForceFullCalc: &t})

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

func SaveExcelFile(f *excelize.File) (string, error) {
	// long path workaround for windows
	outputDir, filename := utils.SplitPath(f.Path)
	return SaveAsExcelFile(f, outputDir, strings.TrimSuffix(filename, filepath.Ext(filename)))

}
