package repository

import (
	"weatherTGBot/internal/usecase/repository"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/logger"
)

type TgBotRepository struct {
	Users        repository.UsersRepository
	Stickers     repository.StickersRepository
	Cities       repository.CitiesRepository
	WeatherTypes repository.WeatherTypeRepository
}

func NewBotRepository(db *postgres.DB, logger logger.Logger) *TgBotRepository {
	return &TgBotRepository{
		Users:        NewUsersRepository(db, logger),
		Stickers:     NewStickersRepository(db, logger),
		Cities:       NewCitiesRepository(db, logger),
		WeatherTypes: NewWeatherTypesRepository(db, logger),
	}
}
