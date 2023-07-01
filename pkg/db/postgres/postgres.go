package postgres

import (
	"database/sql"
	"time"
)
import _ "github.com/lib/pq"

type DB struct {
	*sql.DB
}

func NewDBConnection(dsn string) (*DB, error) {
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

	return &DB{db}, nil
}
