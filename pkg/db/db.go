package db

import "weatherTGBot/internal/usecase/repository"

type TgBotRepo interface {
	repository.UsersRepository
	repository.StickersRepository
}
