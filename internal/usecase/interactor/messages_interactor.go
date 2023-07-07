package interactor

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/pkg/logger"
)

type MessagesInteractor interface {
	SendMessage(ChatID int64, text string, params ...interface{}) error
	SendMessageWithKeyboard(ChatID int64, text string, buttons map[string]string) error
	SendMessageWithRemovingKeyboard(ChatID int64, text string) error

	//todo вынести в отдельный интрактор если распухнет
	GetStickersByType(stickerType string) []string
	SendRandomSticker(message *tgbotapi.Message, stickers []string) error
}

type messagesInteractor struct {
	botAPI *tgbotapi.BotAPI
	repo   *repository.TgBotRepository
	logger logger.Logger
}

func NewMessagesInteractor(botAPI *tgbotapi.BotAPI, repo *repository.TgBotRepository, logger logger.Logger) MessagesInteractor {
	return &messagesInteractor{
		botAPI: botAPI,
		repo:   repo,
		logger: logger,
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
func (i *messagesInteractor) initKeyboard(buttons map[string]string) tgbotapi.ReplyKeyboardMarkup {
	if len(buttons) <= 0 {
		i.logger.Error("Empty Keyboard!")
	}
	var Keyboard = tgbotapi.NewReplyKeyboard()

	for key, _ := range buttons {
		Keyboard.Keyboard = append(Keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(key),
		))
	}

	return Keyboard
}

func (i *messagesInteractor) SendMessageWithKeyboard(ChatID int64, text string, buttons map[string]string) error {
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

// todo сделать айдишником код стикера
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
