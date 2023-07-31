package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"weatherTGBot/internal/domain"
	"weatherTGBot/internal/domain/keyboards"
	"weatherTGBot/internal/domain/weather"
)

// Обрабатывает сообщения авторизации в админку
func (h *MessageHandlerImpl) handleAdminAuthorisationMessage(message *tgbotapi.Message) error {
	switch message.Text {
	case os.Getenv("BOT_ADMIN_SECRET"):
		h.logger.Infof("User [%s] enter to admin panel", message.Chat.FirstName)

		err := h.userInteractor.SetUserState(message.Chat.ID, domain.AdminEnterMenuState)
		if err != nil {
			return err
		}

		text := "Добро пожаловать в админку бота! Выберете действие на клавиатуре"
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, text, keyboards.AdminMainMenuChoice)

	case keyboards.ExitButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.UserUnauthorisedState)
		if err != nil {
			return err
		}

		return h.messagesInteractor.SendMessage(message.Chat.ID, "Нажмите /start для того чтобы узнать состояние погоды")

	default:
		h.logger.Warnf("User: [%s] failed enter to admin panel:", message.Chat.FirstName)
		return h.messagesInteractor.SendMessage(message.Chat.ID, "Неверный пароль! Попробуйте еще раз или введите \"выйти\" для выхода")
	}
}

// Обрабатывает главное меню администратора
func (h *MessageHandlerImpl) handleAdminChoice(message *tgbotapi.Message) error {
	switch message.Text {
	case keyboards.SetStickerButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.AdminAddStickerState)
		if err != nil {
			return err
		}
		return h.messagesInteractor.SendMessageWithRemovingKeyboard(message.Chat.ID, "Выберете стикер для состояния погоды")

	case keyboards.ExitButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.UserAuthState)
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

const exitText = "выйти"

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

		err := h.userInteractor.SetUserState(message.Chat.ID, domain.AdminSetStickerCategoryState)
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
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.AdminEnterMenuState)
		if err != nil {
			return err
		}
		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.AdminMainMenuChoice)
	default:
		text := "Вы не выбрали стикер! введите стикер или \"выйти\" для выхода"
		return h.messagesInteractor.SendMessage(message.Chat.ID, text)
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

		err = h.userInteractor.SetUserState(message.Chat.ID, domain.AdminEnterMenuState)
		if err != nil {
			return err
		}

		return h.messagesInteractor.SendMessageWithKeyboard(message.Chat.ID, "Выберете действие", keyboards.AdminMainMenuChoice)
	}

	switch message.Text {

	case keyboards.ExitButton:
		err := h.userInteractor.SetUserState(message.Chat.ID, domain.AdminEnterMenuState)
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
