package weather

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

func SetWeatherTypes(repo *repository.TgBotRepository) {
	weatherTypes, err := repo.WeatherTypes.GetWeatherTypes()
	if err != nil {
		panic(err)
	}

	for _, weatherType := range weatherTypes {
		WeatherTypes[weatherType.Alias] = weatherType.Title
	}
}

func IsWeatherType(weatherType string) bool {
	_, ok := WeatherTypes[weatherType]
	return ok
}
