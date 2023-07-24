package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/domain"
	"weatherTGBot/internal/domain/cities"
	"weatherTGBot/internal/domain/keyboards"
)

// Обрабатывает главное меню пользователя
func (h *MessageHandlerImpl) handleUserMainMenu(message *tgbotapi.Message) error {
	switch message.Text {
	case keyboards.ShowWeatherButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.UserCityChoiceState)
		if err != nil {
			return err
		}
		text := "Выберете город на клавиатуре, чтобы узнать состояние погоды в нем"
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, text, keyboards.GetCustomKeyboard(cities.Cities))

	case keyboards.ExitButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.UserUnauthorisedState)
		if err != nil {
			return err
		}
		//todo сделать красивую менюшку
		return h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, "Вы вышли из бота!")

	default:
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Стартовое приятное сообщение")
	}
}

// Обрабатывает сообщение с выбором города для погоды
func (h *MessageHandlerImpl) handleCityChoiceMessage(message *tgbotapi.Message) error {
	if message.Text == keyboards.ExitButton {
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.UserAuthState)
		if err != nil {
			return err
		}
		return h.messagesInteractor.SendMessageWithKeyboard(
			message.Chat.ID,
			"Выберете действие",
			keyboards.UserMainMenuChoice)
	}

	if _, ok := cities.Cities[message.Text]; !ok {
		return h.messagesInteractor.SendMessage(
			message.Chat.ID,
			"Вы не выбрали город из предложенных!",
		)
	}

	selectedCity := message.Text

	choseCityMsg := fmt.Sprintf("Вы выбрали город %v", selectedCity)
	if err := h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, choseCityMsg); err != nil {
		return err
	}

	//todo подумать над общими ресурсами
	wr, err := h.weatherInteractor.GetWeatherByCity(selectedCity)
	if err != nil {
		return fmt.Errorf("err get weatherr by city:%v", err)
	}

	err = h.weatherInteractor.SendWeather(message, wr)
	if err != nil {
		return err
	}

	err = h.userInteractor.SetUserState(message.Chat.ID, domain.UserAuthState)
	if err != nil {
		return err
	}

	//todo узнать как можно без сообщения вызывать клавиатуру
	return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.UserMainMenuChoice)
}

// Обрабатывает сообщение от неавторизованных пользователей
func (h *MessageHandlerImpl) handleUnauthorisedUserMessage(message *tgbotapi.Message) error {
	text := "Вы не автризованы! Пожалуйста, нажмите комманду /start для авторизации"
	return h.messagesInteractor.SendMessage(message.Chat.ID, text)
}
