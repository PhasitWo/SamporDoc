package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const AppConfigDirName = "SamporDoc"
const SettingFileName = "setting.json"
const DBFileName = "DB.db"

func GetAppConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting UserConfigDir:", err)
		return "", err
	}

	appConfigDir := filepath.Join(configDir, AppConfigDirName)
	err = os.MkdirAll(appConfigDir, 0755)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return "", err
	}

	return appConfigDir, nil
}

func GetAppFilePath(filename string) (string, error) {
	appConfigDir, err := GetAppConfigDir()
	if err != nil {
		fmt.Println("Error getting app config directory:", err)
		return "", err
	}
	return filepath.Join(appConfigDir, filename), nil
}
