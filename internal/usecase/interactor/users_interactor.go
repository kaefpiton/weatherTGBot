package interactor

import (
	"errors"
	"fmt"
	"weatherTGBot/internal/domain"
	"weatherTGBot/internal/infrastructure/repository"
	repository2 "weatherTGBot/internal/usecase/repository"
	"weatherTGBot/pkg/logger"
)

type UsersInteractor interface {
	InsertUser(firstname, lastname string, chatID int64) error
	IsUserExist(chatID int64) bool
	GetUserStateByChatID(chatID int64) string
	SetUserState(chatID int64, state string) error
}

type usersInteractor struct {
	repo                   *repository.TgBotRepository
	usersStateInmemoryRepo repository2.Inmemory
	logger                 logger.Logger
}

func NewUsersInteractor(repo *repository.TgBotRepository, usersStateInmemoryRepo repository2.Inmemory, logger logger.Logger) UsersInteractor {
	return &usersInteractor{
		repo:                   repo,
		usersStateInmemoryRepo: usersStateInmemoryRepo,
		logger:                 logger,
	}
}

func (i *usersInteractor) InsertUser(firstname, lastname string, chatID int64) error {
	if !i.repo.Users.IsExist(chatID) {
		i.logger.Info("create new user")
		return i.repo.Users.Create(firstname, lastname, domain.UserAuthState, chatID)
	} else {
		i.logger.Info("user already exist. Update last usage")
		return i.repo.Users.UpdateLastUsage(chatID)
	}
}

func (i *usersInteractor) GetUserStateByChatID(chatID int64) string {
	userStateFromCache := i.getUserStateFromCache(chatID)
	if userStateFromCache != "" {
		return userStateFromCache
	}

	if i.repo.Users.IsExist(chatID) {
		state, err := i.repo.Users.GetUserStateByChatID(chatID)
		if err != nil {

		}
		return state
	} else {
		//todo подумать как зарефакторить
		return domain.UserUnauthorisedState
	}
}

func (i *usersInteractor) SetUserState(chatID int64, state string) error {
	if i.repo.Users.IsExist(chatID) {
		//todo сгружать в бд состояние через какое - то время
		i.usersStateInmemoryRepo.Set(chatID, state)
		return i.repo.Users.SetUserState(chatID, state)
	} else {
		return errors.New("user not exist in app")
	}

}

func (i *usersInteractor) IsUserExist(chatID int64) bool {
	return i.repo.Users.IsExist(chatID)
}

// todo зарефакторить с дженериками
func (i *usersInteractor) getUserStateFromCache(chatID int64) string {
	data := i.usersStateInmemoryRepo.Get(chatID)

	if data != nil {
		return fmt.Sprint(data)
	}

	return ""
}
