package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/domain"
	"weatherTGBot/internal/domain/api"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/internal/usecase/interactor"
	"weatherTGBot/pkg/logger"
)

type MessageHandler interface {
	HandleMessage(message *tgbotapi.Message) error
}

type messageHandler struct {
	//	bot                *Bot
	botApi             *tgbotapi.BotAPI
	messagesInteractor interactor.MessagesInteractor
	weatherInteractor  interactor.WeatherInteractor
	userInteractor     interactor.UsersInteractor
	repo               *repository.TgBotRepository
	logger             logger.Logger
}

// todo уйти от зависимости bot
func NewMessageHandler(
	botApi *tgbotapi.BotAPI,
	messagesInteractor interactor.MessagesInteractor,
	weatherInteractor interactor.WeatherInteractor,
	userInteractor interactor.UsersInteractor,
	repo *repository.TgBotRepository,
	logger logger.Logger) *messageHandler {
	return &messageHandler{
		botApi:             botApi,
		messagesInteractor: messagesInteractor,
		weatherInteractor:  weatherInteractor,
		userInteractor:     userInteractor,
		repo:               repo,
		logger:             logger,
	}
}

// Главный обработчик всех сообщений
func (h *messageHandler) HandleMessage(message *tgbotapi.Message) error {
	h.logger.Infof("[%s] %s", message.From.UserName, message.Text)

	//todo для неавторизованных пользователей
	if h.userInteractor.GetUserStateByChatID(message.Chat.ID) == domain.Unauth_usr_state {
		return h.handleUnauthorisedUserMessage(message)
	}

	if h.userInteractor.GetUserStateByChatID(message.Chat.ID) == domain.Auth_usr_state {
		if _, ok := api.Cities[message.Text]; ok {
			return h.handleCityChoiceMessage(message)
		}
	}

	return h.handleDefaultMessage(message)
}

// Обрабатывает сообщение с городом
func (h *messageHandler) handleCityChoiceMessage(message *tgbotapi.Message) error {
	selectedCity := message.Text

	choseCityMsg := fmt.Sprintf("Вы выбрали город %v", selectedCity)
	if err := h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, choseCityMsg); err != nil {
		return err
	}

	//todo подумать над общими ресурсами
	wr, err := h.weatherInteractor.GetWeatherByCity(selectedCity)
	if err != nil {
		return fmt.Errorf("err get weather by city:%v", err)
	}

	return h.weatherInteractor.SendWeather(message, wr)
}

// Обрабатывает сообщение от неавторизованных пользователей
func (h *messageHandler) handleUnauthorisedUserMessage(message *tgbotapi.Message) error {
	return h.messagesInteractor.SendMessage(message.Chat.ID, "Вы не автризовались! Пожалуйста, нажмите комманду start")
}

// Обрабатывает сообщение по умолчанию
func (h *messageHandler) handleDefaultMessage(message *tgbotapi.Message) error {
	return h.messagesInteractor.SendMessage(message.Chat.ID, "Вы не выбрали город ! Пожалуйста, выберете город на клавиатуре")
}
