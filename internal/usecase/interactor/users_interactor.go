package interactor

import (
	"errors"
	"weatherTGBot/internal/domain"
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/pkg/logger"
)

type UsersInteractor interface {
	InsertUser(firstname, lastname string, chatID int64) error
	GetUserStateByChatID(chatID int64) string
	SetUserState(chatID int64, state string) error
}

type usersInteractor struct {
	repo   *repository.TgBotRepository
	logger logger.Logger
}

func NewUsersInteractor(repo *repository.TgBotRepository, logger logger.Logger) UsersInteractor {
	return &usersInteractor{
		repo:   repo,
		logger: logger,
	}
}

func (i *usersInteractor) InsertUser(firstname, lastname string, chatID int64) error {
	if !i.repo.Users.IsExist(chatID) {
		i.logger.Info("create new user")
		return i.repo.Users.Create(firstname, lastname, domain.Auth_usr_state, chatID)
	} else {
		i.logger.Info("user already exist. Update last usage")
		return i.repo.Users.UpdateLastUsage(chatID)
	}
}

// todo прикрутить еще кеш
func (i *usersInteractor) GetUserStateByChatID(chatID int64) string {
	if i.repo.Users.IsExist(chatID) {
		state, err := i.repo.Users.GetUserStateByChatID(chatID)
		if err != nil {

		}
		return state
	} else {
		//todo подумать как зарефакторить
		return domain.Unauth_usr_state
	}

}

// todo прикрутить еще кеш
func (i *usersInteractor) SetUserState(chatID int64, state string) error {
	if i.repo.Users.IsExist(chatID) {
		return i.repo.Users.SetUserState(chatID, state)
	} else {
		//todo refactor err name
		return errors.New("some error")
	}

}
