package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/pkg/logger"
)

type commandHandler interface {
	//todo можно завернуть message чтобы уйти от зависимости api
	handleCommand(message *tgbotapi.Message) error
}

type commandHandlerImpl struct {
	//todo завернуть в интерфейс
	botApi *tgbotapi.BotAPI
	repo   *repository.TgBotRepository
	log    logger.Logger
}

func newCommandHandlerImpl(botApi *tgbotapi.BotAPI, repo *repository.TgBotRepository, log logger.Logger) commandHandler {
	return &commandHandlerImpl{
		botApi: botApi,
		repo:   repo,
		log:    log,
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

// Инициализирует клавиатуру с городами
// todo подумать как можно запилить клавиатуру в 2 строки
func initCitiesKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var Keyboard = tgbotapi.NewReplyKeyboard()

	for key, _ := range cities {
		Keyboard.Keyboard = append(Keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(key),
		))
	}

	return Keyboard
}

// Обрабатывает команду /start
func (h *commandHandlerImpl) handleStartCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle start command:%s", message)

	greetings := "Добро пожаловать " + message.From.FirstName + "!"

	msg := tgbotapi.NewMessage(message.Chat.ID, greetings)
	_, err := h.botApi.Send(msg)
	if err != nil {
		return err
	}

	//todo через интерактор
	h.repo.Users.InsertUser(message.From.FirstName, message.From.LastName, message.Chat.ID)

	msg = tgbotapi.NewMessage(message.Chat.ID, "Выберете город на клавиатуре, чтобы узнать состояние погоды в нем")
	msg.ReplyMarkup = initCitiesKeyboard()
	_, err = h.botApi.Send(msg)

	return err
}

// Обрабатывает команду /info
func (h *commandHandlerImpl) handleInfoCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle info command:%s", message)

	text := "Бот, отсылающий состояние погоды на текущий момент в разных городах России"

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	_, err := h.botApi.Send(msg)

	return err
}

// Обрабатывает отсутствие известной команды
func (h *commandHandlerImpl) handleDefaultCommand(message *tgbotapi.Message) error {
	h.log.Infof("Handle default command:%s", message)

	defaultText := "Я не знаю такой команды :("
	msg := tgbotapi.NewMessage(message.Chat.ID, defaultText)
	_, err := h.botApi.Send(msg)

	return err
}
