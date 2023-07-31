package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/infrastructure/handlers"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/internal/usecase/interactor"
	repository2 "weatherTGBot/internal/usecase/repository"
	"weatherTGBot/pkg/logger"
	"weatherTGBot/pkg/weather"
)

type Bot struct {
	bot            *tgbotapi.BotAPI
	commandHandler handlers.CommandHandler
	messageHandler handlers.MessageHandler
	repo           *repository.TgBotRepository
	logger         logger.Logger
}

func NewBot(botApi *tgbotapi.BotAPI,
	weatherApi weather.WeatherApi,
	repo *repository.TgBotRepository,
	ttl repository2.TLLRepository,
	logger logger.Logger) *Bot {
	tgbot := &Bot{
		bot:    botApi,
		repo:   repo,
		logger: logger,
	}

	messagesInteractor := interactor.NewMessagesInteractor(botApi, repo, logger)
	weatherInteractor := interactor.NewWeatherInteractor(weatherApi, messagesInteractor, ttl)
	usersStateInmemoryRepo := repository.NewUsersStatesInmemoryRepository()
	usersInteractor := interactor.NewUsersInteractor(repo, usersStateInmemoryRepo, logger)
	tgbot.commandHandler = handlers.NewCommandHandler(messagesInteractor, usersInteractor, logger)
	tgbot.messageHandler = handlers.NewMessageHandler(messagesInteractor, weatherInteractor, usersInteractor, logger)

	return tgbot
}

func (b *Bot) Start() error {
	b.logger.Infof("Bot authorised on account: %v", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return nil
}

func (b *Bot) Stop() {
	b.logger.Info("stop receiving updates")
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
			err := b.commandHandler.HandleCommand(update.Message)
			if err != nil {
				b.logger.Error(err)
			}
			continue
		}

		err := b.messageHandler.HandleMessage(update.Message)
		if err != nil {
			b.logger.Error(err)
		}
	}
}
