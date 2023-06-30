package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/pkg/db"
	"weatherTGBot/pkg/logger"
)

type Bot struct {
	bot            *tgbotapi.BotAPI
	commandHandler commandHandler
	weatherApi     WeatherApi
	db             db.TgBotRepo
	log            logger.Logger
}

func NewBot(bot *tgbotapi.BotAPI, weatherApi WeatherApi, db db.TgBotRepo, log logger.Logger) *Bot {
	tgbot := &Bot{
		bot:        bot,
		db:         db,
		log:        log,
		weatherApi: weatherApi,
	}
	tgbot.commandHandler = newCommandHandlerImpl(bot, db, log)

	return tgbot
}

func (b *Bot) Start() error {
	b.log.Info("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) Stop() {
	b.log.Info("stop tg bot")
	b.bot.StopReceivingUpdates()
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		if update.Message.IsCommand() { // обрабатываем команду
			b.commandHandler.handleCommand(update.Message)
			continue
		}
		b.handleMessage(update.Message)
	}
}
