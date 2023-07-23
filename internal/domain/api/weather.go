package api

import (
	"weatherTGBot/internal/infrastructure/repository"
)

type Weather struct {
	Temperature          float64
	TemperatureFeelsLike float64
	Pressure             float64
	WindSpeed            float64
	Clouds               int
	Humidity             int
}

var WeatherTypes = map[string]string{}

func SetWeatherTypes(repo *repository.TgBotRepository) error {
	weatherTypes, err := repo.WeatherTypes.GetWeatherTypes()
	if err != nil {
		return err
	}

	for _, weatherType := range weatherTypes {
		WeatherTypes[weatherType.Alias] = weatherType.Title
	}

	return nil
}

func IsWeatherType(weatherType string) bool {
	_, ok := WeatherTypes[weatherType]
	return ok
}

const HighTemperature = "high temperature"
const NormalTemperature = "normal temperature"
const ColdTemperature = "cold temperature"
const FrostTemperature = "frost temperature"
const PressureHigh = "pressure high"
const PressureNormal = "pressure normal"
const HighWind = "high wind"
const NormalWind = "normal wind"
const LowWind = "low wind"
