package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const jsonConfig = ".gatorconfig.json"

type Config struct {
	DBUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	err := write(*cfg)
	return err
}

func Read() (Config, error) {
	jsonPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

	var cfg Config
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&cfg) ; err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}
	jsonPath := filepath.Join(homeDir, jsonConfig)
	return jsonPath, nil

}

func write(cfg Config) error {
	jsonPath, err := getConfigPath()
	if err != nil {
		return err
	}

	jsonFile, err := os.OpenFile(jsonPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	if err := encoder.Encode(cfg) ; err != nil {
		return err
	}


	return nil
}
