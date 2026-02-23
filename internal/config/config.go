package config

import (
	"encoding/json"
	"os"
	"path"
)

type Config struct {
	DBUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}
	// check file exists
	fullConfigPath := path.Join(home, configFileName)

	// read from file
	data, err := os.ReadFile(fullConfigPath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	// return config
	return config, nil
}

func (c *Config) SetUser(user string) error {
	c.CurrentUser = user
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func write(conf Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	fullConfigPath := path.Join(homeDir, configFileName)
	jsonData, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	err = os.WriteFile(fullConfigPath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
