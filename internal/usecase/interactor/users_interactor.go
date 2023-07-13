package interactor

import (
	"weatherTGBot/internal/infrastructure/repository"
	"weatherTGBot/pkg/logger"
)

type UsersInteractor interface {
	InsertUser(firstname, lastname string, chatID int64) error
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
		return i.repo.Users.Create(firstname, lastname, chatID)
	} else {
		i.logger.Info("user already exist. Update last usage")
		return i.repo.Users.UpdateLastUsage(chatID)
	}
}
