package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const(
	cfgFileName = ".gatorconfig.json"
)


type Config struct {
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath () (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(homeDirectory,cfgFileName)
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			newConf := Config {
				DBURL: "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
				CurrentUserName: "",
			}
			jsonData, err := json.MarshalIndent(newConf, "", "  ")
			if err != nil {
				return "", err
			}
			err = os.WriteFile(filePath, jsonData, 0644)
			if err != nil {
				return "", err
			}
		}
	}
	return filePath, nil
}



func Read() (Config,error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonFile, err := os.Open(cfgPath)
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
    if err != nil {
    	return Config{}, err
    }
	var cfg Config
	err = json.Unmarshal(byteValue, &cfg)
	if err != nil {
    	return Config{}, err
    }
	return cfg,nil
}

func write(cfg Config) error {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return  err
	}
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
        return err
    }
	err = os.WriteFile(cfgPath, jsonData, 0644) // 0644 sets file permissions
    if err != nil {
        return err
    }
	return nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}