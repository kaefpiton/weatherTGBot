package weather

import (
	owm "github.com/briandowns/openweathermap"
	"sync"
)

// todo вынести
type WeatherApi interface {
	GetTemperature() float64
	GetTemperatureFeelsLike() float64
	GetPressure() float64
	GetWindSpeed() float64
	GetCloudPercentage() int
	GetHumidity() int
	SetOptions(options WeatherOptions) error
}

// todo при выносе сократить до options
type WeatherOptions struct {
	City string
}

func NewWeatherOptions(city string) WeatherOptions {
	return WeatherOptions{
		City: city,
	}
}

type openWeatherMapApi struct {
	mu         sync.Mutex
	weatherAPI *owm.CurrentWeatherData
}

func NewOpenWeatherMapApi(unit string, lang string, apiKey string) (WeatherApi, error) {
	weatherAPI, err := owm.NewCurrent(unit, lang, apiKey)
	if err != nil {
		return nil, err
	}
	return &openWeatherMapApi{
		weatherAPI: weatherAPI,
	}, nil
}

func (w *openWeatherMapApi) GetTemperature() float64 {
	return w.weatherAPI.Main.Temp
}

func (w *openWeatherMapApi) GetTemperatureFeelsLike() float64 {
	return w.weatherAPI.Main.FeelsLike
}

func (w *openWeatherMapApi) GetPressure() float64 {
	return convertHpaToMMHG(w.weatherAPI.Main.Pressure)
}

const mmHgInHpa = 1.333

func convertHpaToMMHG(gpa float64) float64 {
	return gpa / mmHgInHpa
}

func (w *openWeatherMapApi) GetWindSpeed() float64 {
	return w.weatherAPI.Wind.Speed
}

func (w *openWeatherMapApi) GetCloudPercentage() int {
	return w.weatherAPI.Clouds.All
}

func (w *openWeatherMapApi) GetHumidity() int {
	return w.weatherAPI.Main.Humidity
}

func (w *openWeatherMapApi) SetOptions(options WeatherOptions) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.weatherAPI.CurrentByName(options.City)
}
