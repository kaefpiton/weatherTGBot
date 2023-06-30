package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Repository struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		DBName   string `json:"db_name"`
		SSLMode  string `json:"ssl_mode"`
	}

	WeatherApi struct {
		APIKey string `json:"weather_api_key"`
		Unit   string `json:"unit"`
		Lang   string `json:"lang"`
	}
	TelegramApi struct {
		APIKey string `json:"telegram_api_key"`
		Debug  bool   `json:"debug"`
	}

	Logger struct {
		Lvl      string `json:"lvl"`
		FilePath string `json:"filePath"`
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

func GetPgDsn(cnf Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cnf.Repository.Host,
		cnf.Repository.Port,
		cnf.Repository.User,
		cnf.Repository.Password,
		cnf.Repository.DBName)
}
