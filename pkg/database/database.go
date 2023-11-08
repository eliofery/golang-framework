package database

import (
	"database/sql"
)

type Database interface {
	Init() (*sql.DB, error)
}

func Init(db Database) (*sql.DB, error) {
	connect, err := db.Init()
	if err != nil {
		return nil, err
	}

	return connect, nil
}
