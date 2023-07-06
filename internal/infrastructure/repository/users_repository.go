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

// todo логика интерактора)
func (r *UsersRepository) InsertUser(firstname, lastname string, chatID int64) error {
	userID, err := r.getUserIDByChatID(r.db, chatID)
	if err != nil {
		return err
	}

	if userID == -1 {
		return r.createUser(r.db, firstname, lastname, chatID)
	} else {
		r.logger.Info("user already exist. Update last usage")
		return r.updateLastUsage(r.db, userID)
	}
}

func (r *UsersRepository) createUser(db *postgres.DB, firstname, lastname string, chatID int64) error {
	r.logger.Info("create new user")

	user := User{
		Firstname: firstname,
		Lastname:  lastname,
		ChatID:    chatID,
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := db.Exec("INSERT INTO users (firstname, lastname, chat_id) values ($1, $2, $3)",
		user.Firstname,
		user.Lastname,
		user.ChatID)

	return err
}

// todo сделать ChatId уникальным полем
func (r *UsersRepository) updateLastUsage(db *postgres.DB, userID int64) error {
	_, err := db.Exec("UPDATE users SET last_usage = $1 WHERE id = $2",
		time.Now(),
		userID)

	return err
}

// todo возможно primary сделать chatID (посмотреть не повторяются ли)
func (r *UsersRepository) getUserIDByChatID(db *postgres.DB, chatID int64) (int64, error) {
	stmt, err := db.Prepare("SELECT ID FROM users WHERE chat_id = $1")
	if err != nil {
		return -1, err
	}

	var userID int64
	r.mu.RLock()
	defer r.mu.Unlock()
	err = stmt.QueryRow(chatID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}

		return -1, err
	}
	return userID, nil
}
