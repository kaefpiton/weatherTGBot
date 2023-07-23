package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"weatherTGBot/internal/domain"
	"weatherTGBot/internal/domain/cities"
	"weatherTGBot/internal/domain/keyboards"
	"weatherTGBot/internal/domain/weather"
	"weatherTGBot/internal/usecase/interactor"
	"weatherTGBot/pkg/logger"
)

type MessageHandler interface {
	HandleMessage(message *tgbotapi.Message) error
}

type MessageHandlerImpl struct {
	messagesInteractor interactor.MessagesInteractor
	weatherInteractor  interactor.WeatherInteractor
	userInteractor     interactor.UsersInteractor
	logger             logger.Logger
}

func NewMessageHandler(
	messagesInteractor interactor.MessagesInteractor,
	weatherInteractor interactor.WeatherInteractor,
	userInteractor interactor.UsersInteractor,
	logger logger.Logger) *MessageHandlerImpl {
	return &MessageHandlerImpl{
		messagesInteractor: messagesInteractor,
		weatherInteractor:  weatherInteractor,
		userInteractor:     userInteractor,
		logger:             logger,
	}
}

// Главный обработчик всех сообщений
func (h *MessageHandlerImpl) HandleMessage(message *tgbotapi.Message) error {
	h.logger.Infof("[%s] %s", message.From.UserName, message.Text)

	state := h.userInteractor.GetUserStateByChatID(message.Chat.ID)
	switch state {
	//Users handlers
	case domain.User_unauth_state:
		return h.handleUnauthorisedUserMessage(message)

	case domain.User_auth_state:
		return h.handleUserChoice(message)

	case domain.User_city_choice_state:
		return h.handleCityChoiceMessage(message)

	//Admins handlers
	case domain.Admin_auth_state:
		return h.handleAdminAuthorisationMessage(message)

	case domain.Admin_enter_state:
		return h.handleAdminChoice(message)

	case domain.Admin_add_sticker_state:
		return h.handleAdminSetSticker(message)

	case domain.Admin_set_sticker_category_state:
		return h.handleAdminSetStickerCategory(message)

	default:
		return h.handleDefaultMessage(message)
	}

}

// Обрабатывает главное меню пользователя
func (h *MessageHandlerImpl) handleUserChoice(message *tgbotapi.Message) error {
	switch message.Text {
	case keyboards.ShowWeatherButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.User_city_choice_state)
		if err != nil {
			return err
		}
		text := "Выберете город на клавиатуре, чтобы узнать состояние погоды в нем"
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, text, keyboards.GetCustomKeyboard(cities.Cities))

	case keyboards.ExitButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.User_unauth_state)
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
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.User_auth_state)
		if err != nil {
			return err
		}
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.UserMainMenuChoice)
	}

	if _, ok := cities.Cities[message.Text]; !ok {
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Вы не выбрали город из предложенных!")
	}

	selectedCity := message.Text

	choseCityMsg := fmt.Sprintf("Вы выбрали город %v", selectedCity)
	if err := h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, choseCityMsg); err != nil {
		return err
	}

	//todo подумать над общими ресурсами
	weather, err := h.weatherInteractor.GetWeatherByCity(selectedCity)
	if err != nil {
		return fmt.Errorf("err get weather by city:%v", err)
	}

	err = h.weatherInteractor.SendWeather(message, weather)
	if err != nil {
		return err
	}

	err = h.userInteractor.SetUserState(message.Chat.ID, domain.User_auth_state)
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

// Обрабатывает сообщения авторизации в админку
func (h *MessageHandlerImpl) handleAdminAuthorisationMessage(message *tgbotapi.Message) error {
	switch message.Text {
	case os.Getenv("BOT_ADMIN_SECRET"):
		h.logger.Infof("User [%s] enter to admin panel", message.Chat.FirstName)

		err := h.userInteractor.SetUserState(message.Chat.ID, domain.Admin_enter_state)
		if err != nil {
			return err
		}

		text := "Добро пожаловать в админку бота! Выберете действие на клавиатуре"
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, text, keyboards.AdminMainMenuChoice)

	case keyboards.ExitButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.User_unauth_state)
		if err != nil {
			return err
		}

		return h.messagesInteractor.SendMessage(message.Chat.ID, "Нажмите /start для того чтобы узнать состояние погоды")

	default:
		h.logger.Warnf("User: [%s] failed enter to admin panel:", message.Chat.FirstName)
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Неверный пароль! Попробуйте еще раз или введите exit для выхода")
	}
}

// Обрабатывает главное меню администратора
func (h *MessageHandlerImpl) handleAdminChoice(message *tgbotapi.Message) error {
	switch message.Text {
	case keyboards.SetStickerButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.Admin_add_sticker_state)
		if err != nil {
			return err
		}
		return h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, "Выберете стикер для состояния погоды")

	case keyboards.ExitButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.User_auth_state)
		if err != nil {
			return err
		}

		err = h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, "Вы вышли из админки! Переход в действия пользователя")
		if err != nil {
			return err
		}
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.UserMainMenuChoice)

	default:
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Выберете действие на клавиатуре!")
	}
}

const exitText = "exit"

// Обрабатывает выбор стикера для администратора
func (h *MessageHandlerImpl) handleAdminSetSticker(message *tgbotapi.Message) error {
	if message.Sticker != nil {
		if h.messagesInteractor.IsStickerExist(message.Sticker.FileID) {
			text := "Такой стикер существует в базе, выберете другой стикер"
			err := h.messagesInteractor.SendMessage(message.Chat.ID, text)
			if err != nil {
				return err
			}
		}

		sticker := domain.NewSticker(message.Sticker.SetName, message.Sticker.FileID)
		h.messagesInteractor.StoreStickerCodeByChatId(message.Chat.ID, sticker)

		err := h.userInteractor.SetUserState(message.Chat.ID, domain.Admin_set_sticker_category_state)
		if err != nil {
			return err
		}

		return h.messagesInteractor.SendMessageWithKeyboard(
			message.Chat.ID,
			"Выберете категорию для стикера",
			keyboards.GetCustomKeyboard(weather.WeatherTypes),
		)
	}

	switch message.Text {
	case exitText:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.Admin_enter_state)
		if err != nil {
			return err
		}
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.AdminMainMenuChoice)
	default:
		//todo вынести "введите exit для выхода"
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Вы не выбрали стикер! введите стикер или exit для выхода")
	}
}

// Обрабатывает выбор категрии для стикера администратором
func (h *MessageHandlerImpl) handleAdminSetStickerCategory(message *tgbotapi.Message) error {
	if weather.IsWeatherType(message.Text) {
		sticker := h.messagesInteractor.GetStickerByChatId(message.Chat.ID)
		err := h.messagesInteractor.CreateSticker(sticker, weather.WeatherTypes[message.Text])
		if err != nil {
			return err
		}

		err = h.messagesInteractor.SendMessageWithRemovingKeyboard(
			message.Chat.ID,
			"Стикер добавлен в категорию  "+message.Text+" !",
		)
		if err != nil {
			return err
		}

		err = h.userInteractor.SetUserState(message.Chat.ID, domain.Admin_enter_state)
		if err != nil {
			return err
		}

		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.AdminMainMenuChoice)
	}

	switch message.Text {

	case keyboards.ExitButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.Admin_enter_state)
		if err != nil {
			return err
		}

		err = h.messagesInteractor.SendMessage(message.Chat.ID, "Выход")
		if err != nil {
			return err
		}

		return h.messagesInteractor.SendMessageWithKeyboard(
			message.Chat.ID,
			"Выберете действие",
			keyboards.AdminMainMenuChoice,
		)

	default:
		text := "Вы не выбрали категорию для стикера! введите категорию или exit для выхода из выбора категории для стикера"
		return h.messagesInteractor.SendMessage(message.Chat.ID, text)
	}
}

func (h *MessageHandlerImpl) handleDefaultMessage(message *tgbotapi.Message) error {
	h.logger.Warnf("User with ChatID [%d] throw to default case", message.Chat.ID)
	return nil
}
