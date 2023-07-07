package repository

import (
	"sync"
	"weatherTGBot/internal/usecase/repository"
	"weatherTGBot/pkg/db/postgres"
	"weatherTGBot/pkg/logger"
)

type Sticker struct {
	ID          int
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

// todo сделать айдишником код стикера
func (r *StickersRepository) GetStickersCodesByType(stickerType string) ([]string, error) {
	var stickerCodes []string

	r.mu.RLock()
	query := `SELECT code FROM stickers JOIN sticker_types st on stickers.type_id = st.id where st.title= $1`
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
