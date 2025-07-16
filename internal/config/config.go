package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {

	configFile, err := getConfigFilePath()

	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(configFile)

	if err != nil {
		return Config{}, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	cfg := Config{}

	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}

func getConfigFilePath() (string, error) {
	homeLocation, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	configFile := filepath.Join(homeLocation, configFileName)

	return configFile, nil
}

func write(cfg Config) error {
	configFile, err := getConfigFilePath()

	if err != nil {
		return err
	}

	file, err := os.Create(configFile)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(cfg); err != nil {
		return err
	}

	return nil
}
