package repository

import (
	"sync"
	"weatherTGBot/internal/usecase/repository"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/logger"
)

type WeatherTypesRepository struct {
	db     *postgres.DB
	mu     sync.RWMutex
	logger logger.Logger
}

func NewWeatherTypesRepository(db *postgres.DB, log logger.Logger) repository.WeatherTypeRepository {
	return &WeatherTypesRepository{
		db:     db,
		logger: log,
	}
}

func (r *WeatherTypesRepository) GetWeatherTypes() ([]repository.WeatherType, error) {
	r.logger.Info("Get weather types from db")

	r.mu.RLock()
	defer r.mu.RUnlock()

	var weatherTypes []repository.WeatherType

	query := "SELECT title,alias FROM weather_types"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var weatherType repository.WeatherType
		err = rows.Scan(&weatherType.Title, &weatherType.Alias)
		if err != nil {
			r.logger.Error(err)
			continue
		}
		weatherTypes = append(weatherTypes, weatherType)
	}

	return weatherTypes, nil
}
