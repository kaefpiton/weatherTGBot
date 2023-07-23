package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/domain"
	"weatherTGBot/internal/domain/keyboards"
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
const commandInfo = "info"
const commandAdminPanel = "admin"

// Главный обработчик всех команд
func (h *commandHandler) HandleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return h.handleStartCommand(message)

	case commandAdminPanel:
		return h.handleAdminPanelCommand(message)

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

	err = h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.UserMainMenuChoice)
	if err != nil {
		return err
	}

	h.usersInteractor.SetUserState(message.Chat.ID, domain.UserAuthState)

	return err
}

// Обрабатывает команду /info
func (h *commandHandler) handleInfoCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle info command:%s", message)

	text := "Бот, отсылающий состояние погоды на текущий момент в разных городах России"

	return h.messagesInteractor.SendMessage(message.Chat.ID, text)
}

// Обрабатывает команду /admin
func (h *commandHandler) handleAdminPanelCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle admin panel command:%s", message)

	if !h.usersInteractor.IsUserExist(message.Chat.ID) {
		text := "Вы не автризовались! Пожалуйста, нажмите комманду /start для авторизации"
		return h.messagesInteractor.SendMessage(message.Chat.ID, text)
	}
	h.usersInteractor.SetUserState(message.Chat.ID, domain.AdminAuthState)

	text := "Введите пароль для админки:"
	return h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, text)
}

// Обрабатывает отсутствие известной команды
func (h *commandHandler) handleDefaultCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle default command:%s", message)

	defaultText := "Я не знаю такой команды :("

	return h.messagesInteractor.SendMessage(message.Chat.ID, defaultText)
}
