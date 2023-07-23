package api

import (
	"weatherTGBot/internal/infrastructure/repository"
)

var Cities = map[string]string{}

func SetCities(repo *repository.TgBotRepository) {
	cities, err := repo.Cities.GetCites()
	if err != nil {
		panic(err)
	}

	for _, city := range cities {
		Cities[city.Alias] = city.Title
	}
}
