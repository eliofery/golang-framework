package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

type Storage struct {
	Path string
}

func New() *Storage {
	return &Storage{
		Path: os.Getenv("SQLITE_PATH"),
	}
}

func (d *Storage) Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", d.Path)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`
    CREATE TABLE IF NOT EXISTS user(
        id INTEGER PRIMARY KEY,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL);
    `)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return db, nil
}
