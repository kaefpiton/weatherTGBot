package repository

import (
	"database/sql"
	"sync"
	"weatherTGBot/internal/usecase/repository"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/logger"
)

type Sticker struct {
	Title       string
	Code        string
	StickerType string
}

type StickersRepository struct {
	db     *postgres.DB
	mu     sync.RWMutex
	logger logger.Logger
}

func NewStickersRepository(db *postgres.DB, log logger.Logger) repository.StickersRepository {
	return &StickersRepository{
		db:     db,
		logger: log,
	}
}

func (r *StickersRepository) GetStickersCodesByType(stickerType string) ([]string, error) {
	var stickerCodes []string

	r.mu.RLock()
	query := "SELECT code FROM stickers JOIN sticker_types type on stickers.type_id = type.id where type.title= $1"
	rows, err := r.db.Query(query, stickerType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	defer r.mu.RUnlock()

	for rows.Next() {
		var stickerCode string
		err = rows.Scan(&stickerCode)
		if err != nil {
			r.logger.Error(err)
			continue
		}
		stickerCodes = append(stickerCodes, stickerCode)
	}

	return stickerCodes, nil
}

func (r *StickersRepository) IsStickerExist(stickerCode string) bool {
	var exists bool
	r.mu.RLock()
	defer r.mu.RUnlock()
	err := r.db.QueryRow("SELECT EXISTS (SELECT code FROM stickers WHERE code=$1)", stickerCode).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		r.logger.Error("error checking if row exists: %v", err)
	}

	return exists
}

func (r *StickersRepository) GetStickerTypes() []string {
	stickerTypes := make([]string, 0, 0)

	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, err := r.db.Query("SELECT title FROM sticker_types")
	if err != nil && err != sql.ErrNoRows {
		r.logger.Error("error checking if row exists: %v", err)
	}

	for rows.Next() {
		var stickerType string
		err = rows.Scan(&stickerType)
		if err != nil {
			r.logger.Error(err)
			continue
		}
		stickerTypes = append(stickerTypes, stickerType)
	}

	return stickerTypes
}

func (r *StickersRepository) CreateSticker(title, code, categoryTitle string) error {
	r.logger.Info("create new sticker")

	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec("INSERT INTO stickers (title, code, type_id) select $1,$2,id from sticker_types where title = $3;",
		title,
		code,
		categoryTitle)

	r.logger.Errorf("ERR:%v", err)
	return err
}
