package repository

import (
	"database/sql"
	"sync"
	"time"
	"weatherTGBot/internal/domain"
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

func (r *UsersRepository) Create(firstname, lastname, state string, chatID int64) error {
	r.logger.Info("create new user")

	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec("INSERT INTO users (firstname, lastname,state, chat_id) values ($1, $2, $3, $4)",
		firstname,
		lastname,
		state,
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

func (r *UsersRepository) GetUserStateByChatID(chatID int64) (string, error) {
	var state string

	r.mu.Lock()
	defer r.mu.Unlock()
	query := `SELECT state FROM users WHERE chat_id=$1`
	row := r.db.QueryRow(query, chatID)
	err := row.Scan(&state)
	if err != nil {
		//todo может возвращать дефолтный стейт
		return domain.ErrorUserState, err
	}

	return state, nil
}

func (r *UsersRepository) SetUserState(chatID int64, state string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec("UPDATE users SET state = $1 WHERE chat_id = $2",
		state,
		chatID)

	return err
}

func (r *UsersRepository) IsExist(chatID int64) bool {
	var exists bool
	r.mu.RLock()
	defer r.mu.RUnlock()
	err := r.db.QueryRow("SELECT EXISTS (SELECT chat_id FROM users WHERE chat_id=$1)", chatID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		r.logger.Error("error checking if row exists: %v", err)
	}

	return exists
}
