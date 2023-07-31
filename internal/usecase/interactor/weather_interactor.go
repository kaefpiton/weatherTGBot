package interactor

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	weather2 "weatherTGBot/internal/domain/weather"
	repository2 "weatherTGBot/internal/usecase/repository"
	"weatherTGBot/pkg/weather"
)

type WeatherInteractor interface {
	GetWeatherByCity(city string) (*weather2.Weather, error)
	SendWeather(message *tgbotapi.Message, weather *weather2.Weather) error
}

type weatherInteractor struct {
	weatherApi         weather.WeatherApi
	messagesInteractor MessagesInteractor
	ttl                repository2.TLLRepository
}

func NewWeatherInteractor(weatherApi weather.WeatherApi, messagesInteractor MessagesInteractor, ttl repository2.TLLRepository) WeatherInteractor {
	return &weatherInteractor{
		weatherApi:         weatherApi,
		messagesInteractor: messagesInteractor,
		ttl:                ttl,
	}
}

func (i *weatherInteractor) GetWeatherByCity(city string) (*weather2.Weather, error) {
	if w, ok := i.ttl.Get(city); ok {
		cache := w.Value.(*weather2.Weather)
		return cache, nil
	}

	weatherOptions := weather.NewWeatherOptions(city)
	if err := i.weatherApi.SetOptions(weatherOptions); err != nil {
		return nil, err
	}

	w := &weather2.Weather{
		Temperature:          i.weatherApi.GetTemperature(),
		TemperatureFeelsLike: i.weatherApi.GetTemperatureFeelsLike(),
		Pressure:             i.weatherApi.GetPressure(),
		WindSpeed:            i.weatherApi.GetWindSpeed(),
		Clouds:               i.weatherApi.GetCloudPercentage(),
		Humidity:             i.weatherApi.GetHumidity(),
	}
	i.ttl.Put(city, repository2.CreateItem(w))

	return w, nil
}

func (i *weatherInteractor) setOptions(city string) error {
	weatherOptions := weather.NewWeatherOptions(city)
	return i.weatherApi.SetOptions(weatherOptions)
}

// Главный метод. Отсылает все, что есть
func (i *weatherInteractor) SendWeather(message *tgbotapi.Message, weather *weather2.Weather) error {

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
const teperatureMessage = "В городе %s температура %.1f °C. Ощущается как %.1f °C."

func (i *weatherInteractor) sendTemperature(message *tgbotapi.Message, weather *weather2.Weather) error {
	city := message.Text
	temp := weather.Temperature
	feelsLike := weather.TemperatureFeelsLike

	err := i.messagesInteractor.SendMessage(message.Chat.ID, teperatureMessage, city, temp, feelsLike)
	if err != nil {
		return err
	}

	return i.sendTemperatureSticker(message, weather)
}

// Отсылает давление
const presureMessage = "Атмосферное давление %.2f мм ртутного столба"

func (i *weatherInteractor) sendPressure(message *tgbotapi.Message, weather *weather2.Weather) error {
	err := i.messagesInteractor.SendMessage(message.Chat.ID, presureMessage, weather.Pressure)
	if err != nil {
		return err
	}

	return i.sendPressureSticker(message, weather)
}

// Отсылает скорость ветра
const windSpeedMessage = "Скорость ветра  %.2f м/с"

func (i *weatherInteractor) sendWindSpeed(message *tgbotapi.Message, weather *weather2.Weather) error {
	err := i.messagesInteractor.SendMessage(message.Chat.ID, windSpeedMessage, weather.WindSpeed)
	if err != nil {
		return err
	}

	return i.sendWindSpeedSticker(message, weather)
}

// Отсылает дополнительную статистику о погоде
const additionalInfoMessage = "Облачность =  %d%v \nВлажность = %d%v \n"

func (i *weatherInteractor) sendAdditionalInfo(message *tgbotapi.Message, weather *weather2.Weather) error {
	return i.messagesInteractor.SendMessage(message.Chat.ID,
		additionalInfoMessage,
		weather.Clouds, "%",
		weather.Humidity, "%")
}

// Стикеры для температуры
const highTemperatureInCelsius = 27
const normalTemperatureInCelsius = 16
const coldTemperatureInCelsius = 0
const frostTemperatureInCelsius = -15

func (i *weatherInteractor) sendTemperatureSticker(message *tgbotapi.Message, weather *weather2.Weather) error {

	switch {
	case weather.Temperature > highTemperatureInCelsius:
		return i.messagesInteractor.SendRandomSticker(
			message,
			i.messagesInteractor.GetStickersByType(weather2.HighTemperature),
		)

	case weather.Temperature > normalTemperatureInCelsius:
		return i.messagesInteractor.SendRandomSticker(
			message,
			i.messagesInteractor.GetStickersByType(weather2.NormalTemperature),
		)

	case weather.Temperature >= coldTemperatureInCelsius:
		return i.messagesInteractor.SendRandomSticker(
			message,
			i.messagesInteractor.GetStickersByType(weather2.ColdTemperature),
		)

	case weather.Temperature < frostTemperatureInCelsius:
		return i.messagesInteractor.SendRandomSticker(
			message,
			i.messagesInteractor.GetStickersByType(weather2.FrostTemperature),
		)

	default:
		return nil
	}
}

// Стикеры для давления
const normalPressureInMmGg = 760

func (i *weatherInteractor) sendPressureSticker(message *tgbotapi.Message, weather *weather2.Weather) error {
	pressure := weather.Pressure

	if pressure > normalPressureInMmGg {
		return i.messagesInteractor.SendRandomSticker(
			message,
			i.messagesInteractor.GetStickersByType(weather2.PressureHigh),
		)
	} else {
		return i.messagesInteractor.SendRandomSticker(
			message,
			i.messagesInteractor.GetStickersByType(weather2.PressureNormal),
		)
	}
}

// Стикеры скорости для ветра
const highWindSpeedMps = 14
const normalWindSpeedMps = 5
const lowWindSpeedMps = 0

func (i *weatherInteractor) sendWindSpeedSticker(message *tgbotapi.Message, weather *weather2.Weather) error {
	windSpeed := weather.WindSpeed
	switch {
	case windSpeed > highWindSpeedMps:
		return i.messagesInteractor.SendRandomSticker(
			message,
			i.messagesInteractor.GetStickersByType(weather2.HighWind),
		)

	case windSpeed > normalWindSpeedMps:
		return i.messagesInteractor.SendRandomSticker(
			message,
			i.messagesInteractor.GetStickersByType(weather2.NormalWind),
		)

	case windSpeed > lowWindSpeedMps:
		return i.messagesInteractor.SendRandomSticker(
			message,
			i.messagesInteractor.GetStickersByType(weather2.LowWind),
		)

	default:
		return nil
	}
}
