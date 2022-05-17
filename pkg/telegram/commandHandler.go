package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//todo описать остальные команды
const commandStart = "start"
const commandStop = "stop"
const commandInfo = "info"

//todo поменять на более короткое значение константы
const commandChangeCity = "changecity"

//Главный обработчик всех команд
func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		{
			return b.handleStartCommand(message)
		}

	case commandInfo:
		{
			return b.handleInfoCommand(message)
		}

	case commandChangeCity:
		{
			return b.handleChangeCityCommand(message)
		}

	default:
		{
			return b.handleDefaultCommand(message)
		}
	}
}

//Обрабатывает команду /start
func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {

	err := startGreetings(b, message)
	if err != nil {
		return err
	}

	userExist, err := b.db.IsUserExist(message.From.FirstName, message.Chat.ID)
	if err != nil {
		return err
	}

	if userExist {
		fmt.Printf("Пользователь существует\n")
		if err := startWithExistedUser(b, message); err != nil {
			return err
		}
	} else {
		fmt.Printf("Создаем нового пользователя\n")
		if err := startWithNewUser(b, message); err != nil {
			return err
		}
	}
	return nil
}

// Приветствие для '/start'
func startGreetings(b *Bot, message *tgbotapi.Message) error {
	greetings := "Добро пожаловать " + message.From.FirstName + "!"
	msg := tgbotapi.NewMessage(message.Chat.ID, greetings)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// Обработка существующего пользователя
func startWithExistedUser(b *Bot, message *tgbotapi.Message) error {
	if err := b.db.UpdateUserDateOfLastUsage(message.Chat.ID); err != nil {
		return err
	}

	userCity, err := b.db.GetUserCity(message.Chat.ID)
	if err != nil {
		return err
	}

	if err := setAPICity(b, userCity); err != nil {
		return err
	}

	message.Text = userCity
	return b.sendWeather(message)
}

// Обработка нового пользователя
func startWithNewUser(b *Bot, message *tgbotapi.Message) error {
	if err := b.db.CreateUser(message.From.FirstName, message.From.LastName, message.Chat.ID); err != nil {
		return err
	}

	if err := sendCityKeyboardChoice(b, message); err != nil {
		return err
	}
	return nil
}

//Обрабатывает команду /info
func (b *Bot) handleInfoCommand(message *tgbotapi.Message) error {
	//todo придумать нормальный текст
	text := "Бот, отсылающий состояние погоды на текущий момент в разных городах России"

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	_, err := b.bot.Send(msg)
	return err
}

//Обрабатывает команду /changecity
func (b *Bot) handleChangeCityCommand(message *tgbotapi.Message) error {
	if err := sendChangeCityInfo(b, message); err != nil {
		return err
	}

	if err := sendCityKeyboardChoice(b, message); err != nil {
		return err
	}
	return nil
}

//Присылает информацию для изменение города по умолчанию
func sendChangeCityInfo(b *Bot, message *tgbotapi.Message) error {
	text := "Вы собираетесь переназначить город по умолчанию"
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

// Присылает меню городов для пользователя
func sendCityKeyboardChoice(b *Bot, message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберете город на клавиатуре, чтобы узнать состояние погоды в нем")
	msg.ReplyMarkup = initCitiesKeyboard()
	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

//Обрабатывает отсутствие известной команды
func (b *Bot) handleDefaultCommand(message *tgbotapi.Message) error {
	defaultText := "Я не знаю такой команды :("
	msg := tgbotapi.NewMessage(message.Chat.ID, defaultText)
	_, err := b.bot.Send(msg)

	return err
}
