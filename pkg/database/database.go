package database

import (
	"database/sql"
	"github.com/eliofery/golang-image/internal/database"
	"github.com/pressly/goose/v3"
)

const (
	pathMigrations = "migrations"
)

type Database interface {
	Init() (*sql.DB, error)
}

func Init(driver Database) (*sql.DB, error) {
	db, err := driver.Init()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB, dialect goose.Dialect) error {
	if err := goose.SetDialect(string(dialect)); err != nil {
		return err
	}

	if err := goose.Up(db, pathMigrations); err != nil {
		panic(err)
	}

	return nil
}

func MigrateFS(db *sql.DB, dialect goose.Dialect) error {
	goose.SetBaseFS(database.EmbedMigrations)

	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dialect)
}
