package telegram

import (
	"github.com/briandowns/openweathermap"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"weatherTGBot/pkg/db"
)

type Bot struct {
	bot *tgbotapi.BotAPI
	weather *openweathermap.CurrentWeatherData
	db db.Datastore
}

func NewBot(bot *tgbotapi.BotAPI, weather *openweathermap.CurrentWeatherData, db db.Datastore) *Bot {
	return &Bot{
		bot:bot,
		weather:weather,
		db:db,
	}
}


func (b *Bot)Start()error  {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil{
		return err
	}

	b.handleUpdates (updates)

	return nil
}

func (b *Bot)initUpdatesChannel() (tgbotapi.UpdatesChannel, error){
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot)handleUpdates(updates tgbotapi.UpdatesChannel)  {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand(){ // обрабатываем команду
			b.handleCommand(update.Message)
			continue
		}
		b.handleMessage(update.Message)
	}
}







