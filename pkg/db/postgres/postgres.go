package postgres

import (
	"database/sql"
	"time"
	"weatherTGBot/pkg/logger"
)
import _ "github.com/lib/pq"

type DB struct {
	*sql.DB
	log logger.Logger
}

func NewDBConnection(dsn string, log logger.Logger) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(time.Minute * 5)

	return &DB{
		db,
		log,
	}, nil
}
