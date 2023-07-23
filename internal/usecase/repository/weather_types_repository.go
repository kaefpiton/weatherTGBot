package repository

type WeatherType struct {
	Title string
	Alias string
}

type WeatherTypeRepository interface {
	GetWeatherTypes() ([]WeatherType, error)
}
