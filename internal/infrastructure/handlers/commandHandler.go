package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/domain/api"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/internal/usecase/interactor"
	"weatherTGBot/pkg/logger"
)

type CommandHandler interface {
	//todo можно завернуть message чтобы уйти от зависимости api
	HandleCommand(message *tgbotapi.Message) error
}

type commandHandler struct {
	messagesInteractor interactor.MessagesInteractor
	usersInteractor    interactor.UsersInteractor
	repo               *repository.TgBotRepository
	log                logger.Logger
}

func NewCommandHandler(messagesInteractor interactor.MessagesInteractor, usersInteractor interactor.UsersInteractor, log logger.Logger) CommandHandler {
	return &commandHandler{
		messagesInteractor: messagesInteractor,
		usersInteractor:    usersInteractor,
		log:                log,
	}
}

const commandStart = "start"
const commandStop = "stop"
const commandInfo = "info"

//todo вынести в интерфейс (usecase) CommandHandler + реализацию в инфру

// Главный обработчик всех команд
func (h *commandHandler) HandleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return h.handleStartCommand(message)

	case commandInfo:
		return h.handleInfoCommand(message)

	default:
		return h.handleDefaultCommand(message)
	}
}

// Обрабатывает команду /start
func (h *commandHandler) handleStartCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle start command:%s", message)

	greetings := fmt.Sprintf("Добро пожаловать, %v!", message.From.FirstName)
	err := h.messagesInteractor.SendMessage(message.Chat.ID, greetings)
	if err != nil {
		return err
	}

	h.usersInteractor.InsertUser(message.From.FirstName, message.From.LastName, message.Chat.ID)

	text := "Выберете город на клавиатуре, чтобы узнать состояние погоды в нем"
	err = h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, text, api.Cities)
	if err != nil {
		return err
	}

	return nil
}

// Обрабатывает команду /info
func (h *commandHandler) handleInfoCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle info command:%s", message)

	text := "Бот, отсылающий состояние погоды на текущий момент в разных городах России"

	return h.messagesInteractor.SendMessage(message.Chat.ID, text)
}

// Обрабатывает отсутствие известной команды
func (h *commandHandler) handleDefaultCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle default command:%s", message)

	defaultText := "Я не знаю такой команды :("

	return h.messagesInteractor.SendMessage(message.Chat.ID, defaultText)
}