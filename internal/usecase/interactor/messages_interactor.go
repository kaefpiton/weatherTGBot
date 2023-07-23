package interactor

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"weatherTGBot/internal/domain"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/pkg/logger"
)

type MessagesInteractor interface {
	SendMessage(ChatID int64, text string, params ...interface{}) error
	SendMessageWithKeyboard(ChatID int64, text string, buttons []string) error
	SendMessageWithRemovingKeyboard(ChatID int64, text string) error

	//todo вынести в отдельный интрактор
	GetStickersByType(stickerType string) []string
	SendRandomSticker(message *tgbotapi.Message, stickers []string) error
	GetStickersTypes() []string
	CreateSticker(sticker *domain.Sticker, categoryTitle string) error
	StoreStickerCodeByChatId(chatID int64, sticker *domain.Sticker)
	GetStickerByChatId(chatID int64) *domain.Sticker
	IsStickerExist(stickerCode string) bool
}

type messagesInteractor struct {
	botAPI       *tgbotapi.BotAPI
	repo         *repository.TgBotRepository
	logger       logger.Logger
	stickerStore map[int64]*domain.Sticker
}

func NewMessagesInteractor(botAPI *tgbotapi.BotAPI, repo *repository.TgBotRepository, logger logger.Logger) MessagesInteractor {
	return &messagesInteractor{
		botAPI:       botAPI,
		repo:         repo,
		logger:       logger,
		stickerStore: make(map[int64]*domain.Sticker),
	}
}

func (i *messagesInteractor) SendMessage(ChatID int64, text string, params ...interface{}) error {
	if len(params) != 0 {
		text = fmt.Sprintf(text, params...)
	}

	msg := tgbotapi.NewMessage(ChatID, text)
	_, err := i.botAPI.Send(msg)

	return err
}

// Инициализирует клавиатуру с кнопками
// todo подумать как можно запилить клавиатуру в 2 строки
func (i *messagesInteractor) initKeyboard(buttons []string) tgbotapi.ReplyKeyboardMarkup {
	if len(buttons) <= 0 {
		i.logger.Error("Empty Keyboard!")
	}
	var Keyboard = tgbotapi.NewReplyKeyboard()

	for _, val := range buttons {
		Keyboard.Keyboard = append(Keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(val),
		))
	}

	return Keyboard
}

func (i *messagesInteractor) SendMessageWithKeyboard(ChatID int64, text string, buttons []string) error {
	msg := tgbotapi.NewMessage(ChatID, text)
	msg.ReplyMarkup = i.initKeyboard(buttons)
	_, err := i.botAPI.Send(msg)

	return err
}

func (i *messagesInteractor) SendMessageWithRemovingKeyboard(ChatID int64, text string) error {
	msg := tgbotapi.NewMessage(ChatID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	if _, err := i.botAPI.Send(msg); err != nil {
		return err
	}

	return nil
}

func (i *messagesInteractor) GetStickersByType(stickerType string) []string {
	stickers, err := i.repo.Stickers.GetStickersCodesByType(stickerType)
	if err != nil {
		i.logger.Error("StickersRepository not found from repo")
	}

	return stickers
}

func (i *messagesInteractor) SendRandomSticker(message *tgbotapi.Message, stickers []string) error {
	msg := tgbotapi.NewStickerShare(message.Chat.ID, stickers[rand.Intn(len(stickers))])
	_, err := i.botAPI.Send(msg)

	return err
}

func (i *messagesInteractor) GetStickersTypes() []string {
	result := make([]string, 0)

	for _, stickerType := range i.repo.Stickers.GetStickerTypes() {
		result = append(result, stickerType.Title)
	}

	return result
}

func (i *messagesInteractor) CreateSticker(sticker *domain.Sticker, categoryTitle string) error {
	return i.repo.Stickers.CreateSticker(sticker.Name, sticker.Code, categoryTitle)
}

func (i *messagesInteractor) StoreStickerCodeByChatId(chatID int64, sticker *domain.Sticker) {
	i.stickerStore[chatID] = sticker
}

func (i *messagesInteractor) GetStickerByChatId(chatID int64) *domain.Sticker {
	if code, ok := i.stickerStore[chatID]; ok {
		return code
	}
	return nil
}

func (i *messagesInteractor) IsStickerExist(stickerCode string) bool {
	return i.repo.Stickers.IsStickerExist(stickerCode)
}
