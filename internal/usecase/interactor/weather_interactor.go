package interactor

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/domain/api"
	"weatherTGBot/pkg/weather"
)

type WeatherInteractor interface {
	GetWeatherByCity(city string) (*api.Weather, error)
	SendWeather(message *tgbotapi.Message, weather *api.Weather) error
}

type weatherInteractor struct {
	weatherApi         weather.WeatherApi
	messagesInteractor MessagesInteractor
}

func NewWeatherInteractor(weatherApi weather.WeatherApi, messagesInteractor MessagesInteractor) WeatherInteractor {
	return &weatherInteractor{
		weatherApi:         weatherApi,
		messagesInteractor: messagesInteractor,
	}
}

func (i *weatherInteractor) GetWeatherByCity(city string) (*api.Weather, error) {
	weatherOptions := weather.NewWeatherOptions(city)
	if err := i.weatherApi.SetOptions(weatherOptions); err != nil {
		return nil, err
	}

	return &api.Weather{
		Temperature:          i.weatherApi.GetTemperature(),
		TemperatureFeelsLike: i.weatherApi.GetTemperatureFeelsLike(),
		Pressure:             i.weatherApi.GetPressure(),
		WindSpeed:            i.weatherApi.GetWindSpeed(),
		Clouds:               i.weatherApi.GetCloudPercentage(),
		Humidity:             i.weatherApi.GetHumidity(),
	}, nil
}

func (i *weatherInteractor) setOptions(city string) error {
	weatherOptions := weather.NewWeatherOptions(city)
	if err := i.weatherApi.SetOptions(weatherOptions); err != nil {
		return err
	}

	return nil
}

// Главный метод. Отсылает все, что есть
func (i *weatherInteractor) SendWeather(message *tgbotapi.Message, weather *api.Weather) error {

	if err := i.sendTemperature(message, weather); err != nil {
		return err
	}
	if err := i.sendPressure(message, weather); err != nil {
		return err
	}
	if err := i.sendWindSpeed(message, weather); err != nil {
		return err
	}
	if err := i.sendAdditionalInfo(message, weather); err != nil {
		return err
	}

	return nil
}

// Отсылает температуру
func (i *weatherInteractor) sendTemperature(message *tgbotapi.Message, weather *api.Weather) error {
	city := message.Text
	temp := weather.Temperature
	feelsLike := weather.TemperatureFeelsLike

	err := i.messagesInteractor.SendMessage(message.Chat.ID, "В городе %s температура %.1f °C. Ощущается как %.1f °C.", city, temp, feelsLike)
	if err != nil {
		return err
	}

	return i.sendTemperatureSticker(message, weather)
}

// Отсылает давление
func (i *weatherInteractor) sendPressure(message *tgbotapi.Message, weather *api.Weather) error {
	err := i.messagesInteractor.SendMessage(message.Chat.ID, "Атмосферное давление %.2f мм ртутного столба", weather.Pressure)
	if err != nil {
		return err
	}

	return i.sendPressureSticker(message, weather)
}

// Отсылает скорость ветра
func (i *weatherInteractor) sendWindSpeed(message *tgbotapi.Message, weather *api.Weather) error {
	err := i.messagesInteractor.SendMessage(message.Chat.ID, "Скорость ветра  %.2f м/с", weather.WindSpeed)
	if err != nil {
		return err
	}

	return i.sendWindSpeedSticker(message, weather)
}

// Отсылает дополнительную статистику о погоде
func (i *weatherInteractor) sendAdditionalInfo(message *tgbotapi.Message, weather *api.Weather) error {
	return i.messagesInteractor.SendMessage(message.Chat.ID,
		"Облачность =  %d%v \nВлажность = %d%v \n",
		weather.Clouds, "%",
		weather.Humidity, "%")
}

// Стикеры для температуры
func (i *weatherInteractor) sendTemperatureSticker(message *tgbotapi.Message, weather *api.Weather) error {

	switch {
	case weather.Temperature > 27:
		return i.messagesInteractor.SendRandomSticker(message, i.messagesInteractor.GetStickersByType("high temperature"))

	case weather.Temperature > 16:
		return i.messagesInteractor.SendRandomSticker(message, i.messagesInteractor.GetStickersByType("normal temperature"))

	case weather.Temperature >= 0:
		return i.messagesInteractor.SendRandomSticker(message, i.messagesInteractor.GetStickersByType("cold temperature"))

	case weather.Temperature < 0:
		return i.messagesInteractor.SendRandomSticker(message, i.messagesInteractor.GetStickersByType("frost temperature"))

	default:
		return nil
	}
}

// Стикеры для давления
func (i *weatherInteractor) sendPressureSticker(message *tgbotapi.Message, weather *api.Weather) error {
	pressure := weather.Pressure

	if pressure > 760 {
		return i.messagesInteractor.SendRandomSticker(message, i.messagesInteractor.GetStickersByType("pressure high"))
	} else {
		return i.messagesInteractor.SendRandomSticker(message, i.messagesInteractor.GetStickersByType("pressure normal"))
	}
}

// Стикеры скорости для ветра
func (i *weatherInteractor) sendWindSpeedSticker(message *tgbotapi.Message, weather *api.Weather) error {
	windSpeed := weather.WindSpeed

	const highWindSpeed = 14
	const normalWindSpeed = 5
	const lowWindSpeed = 0

	switch {
	case windSpeed > highWindSpeed:
		return i.messagesInteractor.SendRandomSticker(message, i.messagesInteractor.GetStickersByType("high wind"))

	case windSpeed > normalWindSpeed:
		return i.messagesInteractor.SendRandomSticker(message, i.messagesInteractor.GetStickersByType("normal wind"))

	case windSpeed > lowWindSpeed:
		return i.messagesInteractor.SendRandomSticker(message, i.messagesInteractor.GetStickersByType("low wind"))

	default:
		return nil
	}
}
