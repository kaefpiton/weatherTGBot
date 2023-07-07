package repository

import (
	"database/sql"
	"sync"
	"time"
	"weatherTGBot/internal/usecase/repository"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/logger"
)

type User struct {
	ID        int64
	Firstname string
	Lastname  string
	ChatID    int64
}

type UsersRepository struct {
	db     *postgres.DB
	mu     sync.RWMutex
	logger logger.Logger
}

func NewUsersRepository(db *postgres.DB, log logger.Logger) repository.UsersRepository {
	return &UsersRepository{
		db:     db,
		logger: log,
	}
}

// todo логика интерактора) чет тип initUser
func (r *UsersRepository) InsertUser(firstname, lastname string, chatID int64) error {
	if !r.IsUserExist(chatID) {
		r.logger.Info("create new user")
		return r.CreateUser(firstname, lastname, chatID)
	} else {
		r.logger.Info("user already exist. Update last usage")
		return r.UpdateLastUsage(chatID)
	}
}

func (r *UsersRepository) CreateUser(firstname, lastname string, chatID int64) error {
	r.logger.Info("create new user")

	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec("INSERT INTO users (firstname, lastname, chat_id) values ($1, $2, $3)",
		firstname,
		lastname,
		chatID)

	return err
}

func (r *UsersRepository) UpdateLastUsage(chatID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec("UPDATE users SET last_usage = $1 WHERE chat_id = $2",
		time.Now(),
		chatID)

	return err
}

func (r *UsersRepository) IsUserExist(chatID int64) bool {
	var exists bool
	r.mu.RLock()
	defer r.mu.RUnlock()
	err := r.db.QueryRow("SELECT EXISTS (SELECT chat_id FROM users WHERE chat_id=$1)", chatID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		r.logger.Error("error checking if row exists: %v", err)
	}

	return exists
}
