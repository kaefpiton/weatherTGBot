package telegram


import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
)

//Стикеры для температуры
func (b *Bot)sendTemperatureSticker(message *tgbotapi.Message) error {
	temperature := b.weather.Main.Temp

	switch {
	case temperature > 27 :{
		return b.sendRandomSticker(message, b.getStickers("high temperature"))
	}
	case temperature > 16 :{
		return b.sendRandomSticker(message, b.getStickers("normal temperature"))
	}
	case temperature >= 0 :{
		return b.sendRandomSticker(message, b.getStickers("cold temperature"))
	}
	case temperature < 0:{
		return b.sendRandomSticker(message, b.getStickers("frost temperature"))
	}
	default:return nil
	}
}

//Стикеры для давления
func (b *Bot)sendPressureSticker(message *tgbotapi.Message) error {
	pressure :=	convertGpaToMMHG(b.weather.Main.GrndLevel)

	if pressure > 760 {
		return b.sendRandomSticker(message, b.getStickers("pressure high"))
	}else{
		return b.sendRandomSticker(message, b.getStickers("pressure normal"))
	}
}


//Стикеры скорости для ветра
func (b *Bot)sendWindSpeedSticker(message *tgbotapi.Message) error {
	windSpeed := b.weather.Wind.Speed
	const highWindSpeed = 14
	const normalWindSpeed = 5
	const lowWindSpeed = 0

	switch {
	case windSpeed > highWindSpeed:{
		return b.sendRandomSticker(message, b.getStickers("high wind"))
	}
	case windSpeed > normalWindSpeed:{
		return b.sendRandomSticker(message, b.getStickers("normal wind"))
	}
	case windSpeed > lowWindSpeed:{
		return b.sendRandomSticker(message, b.getStickers("low wind"))
	}

	default: return nil
	}
}



func (b Bot)getStickers(stickerType string) ([]string)  {
	stickers, err := b.db.GetStickersCodesByType(stickerType)
	if err != nil{
		log.Printf("Stickers not found from db")
	}
	return stickers
}
func (b *Bot) sendRandomSticker(message *tgbotapi.Message, stickers []string)error{
	msg := tgbotapi.NewStickerShare(message.Chat.ID, stickers[rand.Intn(len(stickers))])
	_, err := b.bot.Send(msg)
	return err
}
