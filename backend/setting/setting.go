package setting

import (
	"SamporDoc/backend/config"
	"SamporDoc/backend/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const DEFAULT_CUSTOMER_DB_PATH = "{{DEFAULT}}"

type Setting struct {
	CustomerDBPath string // "{{DEFAULT}} or custom path"
}

func New() *Setting {
	return &Setting{CustomerDBPath: DEFAULT_CUSTOMER_DB_PATH}
}

func (s *Setting) newError(err error) error {
	return fmt.Errorf("[SETTING]: %w", err)
}

func (s *Setting) getSettingFilePath() (string, error) {
	settingFilePath, err := config.GetAppFilePath(config.SettingFileName)
	if err != nil {
		fmt.Println("Error getting setting file path:", err)
		return "", s.newError(err)
	}
	return settingFilePath, nil
}

func (s *Setting) LoadSetting() error {
	settingFilePath, err := s.getSettingFilePath()
	if err != nil {
		return err
	}
	// check if the file exists
	_, err = os.Stat(settingFilePath)
	if errors.Is(err, os.ErrNotExist) {
		// file does not exist
		file, err := os.Create(settingFilePath)
		if err != nil {
			fmt.Println("Error creating setting.json file:", err)
			return s.newError(err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		err = encoder.Encode(*s)
		if err != nil {
			fmt.Println("Error encoding setting.json file:", err)
			return s.newError(err)
		}
		return nil
	}

	file, err := os.Open(settingFilePath)
	if err != nil {
		fmt.Println("Error opening setting.json file:", err)
		return s.newError(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&s)
	if err != nil {
		fmt.Println("Error decoding setting.json file:", err)
		return s.newError(err)
	}
	return nil
}

func (s *Setting) SaveSetting() error {
	settingFilePath, err := s.getSettingFilePath()
	if err != nil {
		return err
	}
	var file *os.File
	// check if the file exists
	if !utils.FileExists(settingFilePath) {
		// file does not exist
		file, err := os.Create(settingFilePath)
		if err != nil {
			fmt.Println("Error creating setting.json file:", err)
			return s.newError(err)
		}
		defer file.Close()
	}

	file, err = os.OpenFile(settingFilePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening setting.json file:", err)
		return s.newError(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(*s)
	if err != nil {
		fmt.Println("Error encoding setting.json file:", err)
		return s.newError(err)
	}
	return nil
}
