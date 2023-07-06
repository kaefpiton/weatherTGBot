package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/pkg/logger"
)

type Bot struct {
	bot            *tgbotapi.BotAPI
	commandHandler commandHandler
	messageHandler messageHandler
	weatherApi     WeatherApi
	repo           *repository.TgBotRepository
	log            logger.Logger
}

func NewBot(bot *tgbotapi.BotAPI, weatherApi WeatherApi, repo *repository.TgBotRepository, log logger.Logger) *Bot {
	tgbot := &Bot{
		bot:        bot,
		repo:       repo,
		log:        log,
		weatherApi: weatherApi,
	}
	tgbot.commandHandler = newCommandHandlerImpl(bot, repo, log)
	//todo уйти от зввисимости tgbot
	tgbot.messageHandler = newMessageHandlerImpl(tgbot, bot, weatherApi, repo, log)

	return tgbot
}

func (b *Bot) Start() error {
	b.log.Infof("Authorized on account: %v", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) Stop() {
	//todo change log
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
		// ignore any non-Message Updates
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			err := b.commandHandler.handleCommand(update.Message)
			if err != nil {
				b.log.Error(err)
			}
			continue
		}

		err := b.messageHandler.handleMessage(update.Message)
		if err != nil {
			b.log.Error(err)
		}
	}
}
