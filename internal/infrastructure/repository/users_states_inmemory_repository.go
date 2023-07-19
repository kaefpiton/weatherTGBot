package repository

import "weatherTGBot/internal/usecase/repository"

type UsersStatesInmemoryRepository struct {
	cache map[any]any
}

func NewUsersStatesInmemoryRepository() repository.Inmemory {
	return &UsersStatesInmemoryRepository{
		cache: make(map[any]any),
	}
}

func (r *UsersStatesInmemoryRepository) Get(chatId any) any {
	if state, ok := r.cache[chatId]; ok {
		return state
	}

	return nil
}

func (r *UsersStatesInmemoryRepository) Set(chatId any, state any) {
	r.cache[chatId] = state
}
