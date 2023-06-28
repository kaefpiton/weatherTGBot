package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var cities = map[string]string{"Москва": "Moscow", "Ростов": "Rostov", "Агалатово": "Agalatovo"}

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

// Главный обработчик всех сообщений
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	//log.Printf("[%s] %s", message.From.UserName, message.Text)

	if _, ok := cities[message.Text]; ok {
		return b.handleMessageCity(message)
	}
	return b.handleMessageDefault(message)
}

// Обрабатывает сообщение с городом
func (b *Bot) handleMessageCity(message *tgbotapi.Message) error {
	selectedCity := message.Text

	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы выбрали город "+selectedCity)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) //Убирает клавиатуру
	if _, err := b.bot.Send(msg); err != nil {
		return err
	}

	if city, ok := cities[selectedCity]; ok {
		if err := b.setCity(city); err != nil {
			return err
		}
	} else {
		return errors.New("hashmap: The selected city is not in the cities hashmap")
	}

	return b.sendWeather(message)
}

// Обрабатывает сообщение по умолчанию
func (b *Bot) handleMessageDefault(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы не выбрали город ! Пожалуйста, выберете город на клавиатуре")
	_, err := b.bot.Send(msg)
	return err
}

// Отправляет выбранный город API погоды
func (b *Bot) setCity(city string) error {
	return b.weather.CurrentByName(city)
}
