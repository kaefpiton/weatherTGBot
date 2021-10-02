package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//Главный метод. Отсылает все, что есть
func (b *Bot)sendWeather(message *tgbotapi.Message) error {

	if err := b.sendTemperature(message); err != nil {
		return err
	}
	if err := b.sendPressure(message); err != nil {
		return err
	}
	if err := b.sendWindSpeed(message); err != nil {
		return err
	}
	if err := b.sendAdditionalInfo(message); err != nil {
		return err
	}

	return nil
}

//Отсылает температуру
func (b *Bot)sendTemperature(message *tgbotapi.Message) error{
	city := message.Text
	temp := b.weather.Main.Temp
	feelsLike := b.weather.Main.FeelsLike

	msgText := fmt.Sprintf("В городе %s температура %.1f °C. Ощущается как %.1f °C.", city, temp, feelsLike)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err := b.bot.Send(msg)
	if err != nil{
		return err
	}

	return b.sendTemperatureSticker(message)
}

//Отсылает давление
func (b *Bot)sendPressure(message *tgbotapi.Message) error{
	pressure :=	convertGpaToMMHG(b.weather.Main.GrndLevel)
	msgText := fmt.Sprintf("Атмосферное давление %.2f мм ртутного столба", pressure)

	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err := b.bot.Send(msg)
	if err != nil{
		return err
	}

	return b.sendPressureSticker(message)
}

//Отсылает скорость ветра
func (b *Bot)sendWindSpeed(message *tgbotapi.Message) error{
	windSpeed := b.weather.Wind.Speed

	msgText := fmt.Sprintf("Скорость ветра  %.2f м/с", windSpeed)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err := b.bot.Send(msg)
	if err != nil{
		return err
	}

	return b.sendWindSpeedSticker(message)
}

//Отсылает дополнительную статистику о погоде
func (b *Bot)sendAdditionalInfo(message *tgbotapi.Message) error{
	clouds := b.weather.Clouds.All
	humidity := b.weather.Main.Humidity

	msgText := fmt.Sprintf("Облачность =  %d%v \nВлажность = %d%v \n", clouds, "%", humidity, "%")
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	_, err := b.bot.Send(msg)
	return err
}




//todo Вынести функцию в спомогательные
func convertGpaToMMHG(gpa float64) float64{
	return gpa / 1.333
}