package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/domain/api"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/internal/usecase/interactor"
	"weatherTGBot/pkg/logger"
)

type commandHandler interface {
	//todo можно завернуть message чтобы уйти от зависимости api
	handleCommand(message *tgbotapi.Message) error
}

type commandHandlerImpl struct {
	messagesInteractor interactor.MessagesInteractor
	repo               *repository.TgBotRepository
	log                logger.Logger
}

func newCommandHandlerImpl(messagesInteractor interactor.MessagesInteractor, repo *repository.TgBotRepository, log logger.Logger) commandHandler {
	return &commandHandlerImpl{
		messagesInteractor: messagesInteractor,
		repo:               repo,
		log:                log,
	}
}

const commandStart = "start"
const commandStop = "stop"
const commandInfo = "info"

//todo вынести в интерфейс (usecase) commandHandler + реализацию в инфру

// Главный обработчик всех команд
func (h *commandHandlerImpl) handleCommand(message *tgbotapi.Message) error {
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
func (h *commandHandlerImpl) handleStartCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle start command:%s", message)

	greetings := fmt.Sprintf("Добро пожаловать, %v!", message.From.FirstName)
	err := h.messagesInteractor.SendMessage(message.Chat.ID, greetings)
	if err != nil {
		return err
	}

	//todo через интерактор
	h.repo.Users.InsertUser(message.From.FirstName, message.From.LastName, message.Chat.ID)

	text := "Выберете город на клавиатуре, чтобы узнать состояние погоды в нем"
	err = h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, text, api.Cities)
	if err != nil {
		return err
	}

	return nil
}

// Обрабатывает команду /info
func (h *commandHandlerImpl) handleInfoCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle info command:%s", message)

	text := "Бот, отсылающий состояние погоды на текущий момент в разных городах России"

	return h.messagesInteractor.SendMessage(message.Chat.ID, text)
}

// Обрабатывает отсутствие известной команды
func (h *commandHandlerImpl) handleDefaultCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle default command:%s", message)

	defaultText := "Я не знаю такой команды :("

	return h.messagesInteractor.SendMessage(message.Chat.ID, defaultText)
}
