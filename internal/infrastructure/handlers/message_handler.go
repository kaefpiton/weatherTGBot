package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"weatherTGBot/internal/domain"
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

// HandleMessage Главный обработчик всех сообщений
func (h *MessageHandlerImpl) HandleMessage(message *tgbotapi.Message) error {
	h.logger.Infof("[%s] %s", message.From.UserName, message.Text)

	state := h.userInteractor.GetUserStateByChatID(message.Chat.ID)
	switch state {

	//----------Users handlers----------
	case domain.UserUnauthorisedState:
		return h.handleUnauthorisedUserMessage(message)

	case domain.UserAuthState:
		return h.handleUserMainMenu(message)

	case domain.UserCityChoiceState:
		return h.handleCityChoiceMessage(message)

	//----------Admins handlers----------
	case domain.AdminAuthState:
		return h.handleAdminAuthorisationMessage(message)

	case domain.AdminEnterMenuState:
		return h.handleAdminChoice(message)

	case domain.AdminAddStickerState:
		return h.handleAdminSetSticker(message)

	case domain.AdminSetStickerCategoryState:
		return h.handleAdminSetStickerCategory(message)

	//----------Default handler----------
	default:
		return h.handleDefaultMessage(message)
	}

}

func (h *MessageHandlerImpl) handleDefaultMessage(message *tgbotapi.Message) error {
	h.logger.Warnf("User with ChatID [%d] throw to default case", message.Chat.ID)
	return nil
}
