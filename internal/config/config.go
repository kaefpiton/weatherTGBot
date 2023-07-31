package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Postgres struct {
		User         string `json:"user"`
		Password     string `json:"password"`
		Host         string `json:"host"`
		Port         int    `json:"port"`
		DBName       string `json:"db_name"`
		SSLMode      string `json:"ssl_mode"`
		MaxOpenConns int    `json:"maxOpenConns"`
		MaxIdleConns int    `json:"maxIdleConns"`
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

	Weather struct {
		CacheExpireMinutes          int `json:"cacheExpireMinutes"`
		CheckCacheExpireEverySecond int `json:"checkCacheExpireEverySecond"`
	}
}

func LoadConfig(path string) (*Config, error) {
	var config *Config

	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) GetPgDsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.Password,
		c.Postgres.DBName,
		c.Postgres.SSLMode)
}
