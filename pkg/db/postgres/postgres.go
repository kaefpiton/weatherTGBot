package postgres

import "database/sql"
import _"github.com/lib/pq"

type DB struct {
	*sql.DB
}

func NewDB(dataSourceName string)(*DB, error)  {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
