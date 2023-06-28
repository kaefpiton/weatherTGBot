package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const commandStart = "start"
const commandStop = "stop"
const commandInfo = "info"

//todo вынести в интерфейс (usecase) commandHandler + реализацию в инфру

// Главный обработчик всех команд
func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)

	case commandInfo:
		return b.handleInfoCommand(message)

	default:
		return b.handleDefaultCommand(message)
	}
}

// Обрабатывает команду /start
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	b.log.Info("Handle start command:", message)

	greetings := "Добро пожаловать " + message.From.FirstName + "!"

	msg := tgbotapi.NewMessage(message.Chat.ID, greetings)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	//todo через интерактор
	b.db.InsertUser(message.From.FirstName, message.From.LastName, message.Chat.ID)

	msg = tgbotapi.NewMessage(message.Chat.ID, "Выберете город на клавиатуре, чтобы узнать состояние погоды в нем")
	msg.ReplyMarkup = initCitiesKeyboard()
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// Обрабатывает команду /stop
func (b *Bot) handleInfoCommand(message *tgbotapi.Message) error {
	b.log.Info("Handle info command:", message)

	text := "Бот, отсылающий состояние погоды на текущий момент в разных городах России"

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	_, err := b.bot.Send(msg)
	return err

}

// Обрабатывает отсутствие известной команды
func (b *Bot) handleDefaultCommand(message *tgbotapi.Message) error {
	b.log.Info("Handle default command:", message)

	defaultText := "Я не знаю такой команды :("
	msg := tgbotapi.NewMessage(message.Chat.ID, defaultText)
	_, err := b.bot.Send(msg)

	return err
}
