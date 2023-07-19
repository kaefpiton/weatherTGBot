package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"weatherTGBot/internal/domain"
	"weatherTGBot/internal/domain/api"
	"weatherTGBot/internal/domain/keyboards"
	"weatherTGBot/internal/usecase/interactor"
	"weatherTGBot/pkg/logger"
)

type MessageHandler interface {
	HandleMessage(message *tgbotapi.Message) error
}

type messageHandler struct {
	messagesInteractor interactor.MessagesInteractor
	weatherInteractor  interactor.WeatherInteractor
	userInteractor     interactor.UsersInteractor
	logger             logger.Logger
}

func NewMessageHandler(
	messagesInteractor interactor.MessagesInteractor,
	weatherInteractor interactor.WeatherInteractor,
	userInteractor interactor.UsersInteractor,
	logger logger.Logger) *messageHandler {
	return &messageHandler{
		messagesInteractor: messagesInteractor,
		weatherInteractor:  weatherInteractor,
		userInteractor:     userInteractor,
		logger:             logger,
	}
}

// Главный обработчик всех сообщений
func (h *messageHandler) HandleMessage(message *tgbotapi.Message) error {
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

	case domain.Admin_set_sticker_state:
		return h.handleAdminSetSticker(message)

	case domain.Admin_set_sticker_category:
		return h.handleAdminSetStickerCategory(message)

	default:
		return h.handleDefaultMessage(message)
	}

}

// Обрабатывает
func (h *messageHandler) handleUserChoice(message *tgbotapi.Message) error {
	switch message.Text {
	case keyboards.ShowWeatherButton:
		h.userInteractor.SetUserState(message.Chat.ID, domain.User_city_choice_state)
		text := "Выберете город на клавиатуре, чтобы узнать состояние погоды в нем"
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, text, api.GetCitiesKeys(api.Cities))

	case keyboards.ExitButton:
		h.userInteractor.SetUserState(message.Chat.ID, domain.User_unauth_state)
		//todo сделать красивую менюшку
		return h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, "Вы вышли из бота!")

	default:
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Стартовое приятное сообщение")
	}
}

// Обрабатывает сообщение с городом
func (h *messageHandler) handleCityChoiceMessage(message *tgbotapi.Message) error {
	if _, ok := api.Cities[message.Text]; !ok {
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

	h.weatherInteractor.SendWeather(message, weather)

	h.userInteractor.SetUserState(message.Chat.ID, domain.User_auth_state)

	//todo узнать как можно без сообщения вызывать клавиатуру
	return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.UserMainMenuChoice)
}

// Обрабатывает сообщение от неавторизованных пользователей
func (h *messageHandler) handleUnauthorisedUserMessage(message *tgbotapi.Message) error {
	text := "Вы не автризованы! Пожалуйста, нажмите комманду /start для авторизации"
	return h.messagesInteractor.SendMessage(message.Chat.ID, text)
}

// Обрабатывает сообщения авторизации в админку
func (h *messageHandler) handleAdminAuthorisationMessage(message *tgbotapi.Message) error {
	switch message.Text {
	case os.Getenv("BOT_ADMIN_SECRET"):
		h.logger.Infof("User [%s] enter to admin panel", message.Chat.FirstName)

		h.userInteractor.SetUserState(message.Chat.ID, domain.Admin_enter_state)

		text := "Добро пожаловать в админку бота! Выберете действие на клавиатуре"
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, text, keyboards.AdminMainMenuChoice)

	case keyboards.ExitButton:
		h.userInteractor.SetUserState(message.Chat.ID, domain.User_unauth_state)
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Нажмите /start для того чтобы узнать состояние погоды")

	default:
		h.logger.Warnf("User: [%s] failed enter to admin panel:", message.Chat.FirstName)
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Неверный пароль! Попробуйте еще раз или введите exit для выхода")
	}
}

func (h *messageHandler) handleAdminChoice(message *tgbotapi.Message) error {
	switch message.Text {
	case keyboards.SetStickerButton:
		h.userInteractor.SetUserState(message.Chat.ID, domain.Admin_set_sticker_state)
		return h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, "Выберете стикер для состояния погоды")

	case keyboards.ExitButton:
		h.userInteractor.SetUserState(message.Chat.ID, domain.User_auth_state)
		h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, "Вы вышли из админки! Переход в действия пользователя")
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.UserMainMenuChoice)

	default:
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Выберете действие на клавиатуре!")
	}
}

const exitText = "exit"

func (h *messageHandler) handleAdminSetSticker(message *tgbotapi.Message) error {
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

		err := h.userInteractor.SetUserState(message.Chat.ID, domain.Admin_set_sticker_category)
		if err != nil {
			return err
		}

		//todo мб custom keyboard или чет такое
		buttons := h.messagesInteractor.GetStickersTypes()
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете категорию для стикера", buttons)
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

func (h *messageHandler) handleAdminSetStickerCategory(message *tgbotapi.Message) error {
	StickerTypeButtons := h.messagesInteractor.GetStickersTypes()

	if IsStickerType(message.Text, StickerTypeButtons) {
		sticker := h.messagesInteractor.GetStickerByChatId(message.Chat.ID)
		err := h.messagesInteractor.CreateSticker(sticker, message.Text)
		if err != nil {
			return err
		}

		err = h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, "Стикер добавлен в категорию  "+message.Text)
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

		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.AdminMainMenuChoice)

	default:
		text := "Вы не выбрали категорию для стикера! введите категорию или exit для выхода из выбора категории для стикера"
		return h.messagesInteractor.SendMessage(message.Chat.ID, text)
	}
}

// todo refactor
func IsStickerType(stickerType string, stickerTypes []string) bool {
	for _, stype := range stickerTypes {
		if stickerType == stype {
			return true
		}
	}

	return false
}

func (h *messageHandler) handleAdminExitFromAddStickerMenu(message *tgbotapi.Message) error {
	switch message.Text {
	case keyboards.ExitToAdminPanelButton:
		return h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, "Переходим в панель админа")
	case keyboards.AddStickerButton:
		return h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, "Создаем новый стикер")
	default:
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Вы не выбрали действие на клавиатуре")
	}
}

func (h *messageHandler) handleDefaultMessage(message *tgbotapi.Message) error {
	h.logger.Warnf("User with ChatID [%d] throw to default case", message.Chat.ID)
	return nil
}
