package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/pkg/db"
	"weatherTGBot/pkg/logger"
)

type messageHandler interface {
	handleMessage(message *tgbotapi.Message) error
}

type messageHandlerImpl struct {
	bot        *Bot
	botApi     *tgbotapi.BotAPI
	weatherApi WeatherApi
	repo       db.TgBotRepo
	log        logger.Logger
}

// todo уйти от зависимости bot
func newMessageHandlerImpl(bot *Bot, botApi *tgbotapi.BotAPI, weatherApi WeatherApi, repo db.TgBotRepo, log logger.Logger) *messageHandlerImpl {
	return &messageHandlerImpl{
		bot:        bot,
		botApi:     botApi,
		weatherApi: weatherApi,
		repo:       repo,
		log:        log,
	}
}

// todo вынемти
var cities = map[string]string{"Москва": "Moscow", "Ростов": "Rostov", "Агалатово": "Agalatovo"}

// Главный обработчик всех сообщений
func (h *messageHandlerImpl) handleMessage(message *tgbotapi.Message) error {
	//log.Printf("[%s] %s", message.From.UserName, message.Text)

	if _, ok := cities[message.Text]; ok {
		return h.handleCityMessage(message)
	}

	return h.handleDefaultMessage(message)
}

// Обрабатывает сообщение с городом
func (h *messageHandlerImpl) handleCityMessage(message *tgbotapi.Message) error {
	selectedCity := message.Text

	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы выбрали город "+selectedCity)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true) //Убирает клавиатуру
	if _, err := h.botApi.Send(msg); err != nil {
		return err
	}

	if city, ok := cities[selectedCity]; ok {
		weatherOptions := NewWeatherOptions(city)
		if err := h.weatherApi.SetOptions(weatherOptions); err != nil {
			return err
		}
	} else {
		return errors.New("hashmap: The selected city is not in the cities hashmap")
	}

	return h.bot.sendWeather(message)
}

// Обрабатывает сообщение по умолчанию
func (h *messageHandlerImpl) handleDefaultMessage(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Вы не выбрали город ! Пожалуйста, выберете город на клавиатуре")
	_, err := h.botApi.Send(msg)
	return err
}
