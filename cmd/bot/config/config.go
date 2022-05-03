package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DataBase struct {
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"db_name"`
		SSLMode  string `json:"ssl_mode"`
	}

	Weather struct {
		APIKey string `json:"weather_api_key"`
		Unit   string `json:"unit"`
		Lang   string `json:"lang"`
	}
	Telegram struct {
		APIKey string `json:"telegram_api_key"`
		Debug  bool   `json:"debug"`
	}
}

func LoadConfiguration(file string) Config {
	var config Config

	configFile, err := os.Open(file)

	if err != nil {
		fmt.Println(err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}
