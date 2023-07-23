package repository

import (
	"sync"
	"weatherTGBot/internal/usecase/repository"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/logger"
)

type CitiesRepository struct {
	db     *postgres.DB
	mu     sync.RWMutex
	logger logger.Logger
}

func NewCitiesRepository(db *postgres.DB, log logger.Logger) repository.CitiesRepository {
	return &CitiesRepository{
		db:     db,
		logger: log,
	}
}

func (r *CitiesRepository) GetCites() ([]repository.City, error) {
	r.logger.Info("Get cities from db")

	r.mu.RLock()
	defer r.mu.RUnlock()

	var cities []repository.City
	query := "SELECT title,alias FROM cities"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var city repository.City
		err = rows.Scan(&city.Title, &city.Alias)
		if err != nil {
			r.logger.Error(err)
			continue
		}
		cities = append(cities, city)
	}

	return cities, nil
}
