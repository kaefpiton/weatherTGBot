package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/domain/api"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/internal/usecase/interactor"
	"weatherTGBot/pkg/logger"
)

type messageHandler interface {
	handleMessage(message *tgbotapi.Message) error
}

type messageHandlerImpl struct {
	//	bot                *Bot
	botApi             *tgbotapi.BotAPI
	messagesInteractor interactor.MessagesInteractor
	weatherInteractor  interactor.WeatherInteractor
	repo               *repository.TgBotRepository
	logger             logger.Logger
}

// todo уйти от зависимости bot
func newMessageHandlerImpl(
	botApi *tgbotapi.BotAPI,
	messagesInteractor interactor.MessagesInteractor,
	weatherInteractor interactor.WeatherInteractor,
	repo *repository.TgBotRepository,
	logger logger.Logger) *messageHandlerImpl {
	return &messageHandlerImpl{
		botApi:             botApi,
		messagesInteractor: messagesInteractor,
		weatherInteractor:  weatherInteractor,
		repo:               repo,
		logger:             logger,
	}
}

// Главный обработчик всех сообщений
func (h *messageHandlerImpl) handleMessage(message *tgbotapi.Message) error {
	h.logger.Infof("[%s] %s", message.From.UserName, message.Text)

	if _, ok := api.Cities[message.Text]; ok {
		return h.handleCityChoiceMessage(message)
	}

	return h.handleDefaultMessage(message)
}

// Обрабатывает сообщение с городом
func (h *messageHandlerImpl) handleCityChoiceMessage(message *tgbotapi.Message) error {
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

// Обрабатывает сообщение по умолчанию
func (h *messageHandlerImpl) handleDefaultMessage(message *tgbotapi.Message) error {
	return h.messagesInteractor.SendMessage(message.Chat.ID, "Вы не выбрали город ! Пожалуйста, выберете город на клавиатуре")
}
