package main

import (
	"database/sql"
	"github.com/eliofery/golang-image/internal/app/http/controllers/home"
	mw "github.com/eliofery/golang-image/internal/app/http/middleware"
	"github.com/eliofery/golang-image/pkg/config"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/database/postgres"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"os"
)

func main() {
	// Подключение конфигурационного файла .env
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	// Подключение логирования
	logger := logging.New()

	// Подключение к БД Postgres
	db, err := database.Init(postgres.New())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(db)

	// Миграция БД Postgres
	if err = database.MigrateFS(db, goose.DialectPostgres); err != nil {
		logger.Error(err.Error())
	}

	// Создание роутера
	route := router.New()

	// Пользовательский Middleware
	route.Use(middleware.RequestID)
	route.Use(middleware.RealIP)
	route.Use(middleware.Logger)
	route.Use(middleware.Recoverer)
	route.Use(middleware.URLFormat)
	route.Use(mw.Bind(logger, db))
	route.Use(mw.Csrf)

	// Подключение ресурсов
	tpl.AssetsFsInit(route.Mux)

	// Роуты
	route.Get("/", home.Index)
	route.Post("/", home.Post)

	// Запуск сервера
	logger.Info("Сервер запущен: http://localhost:8080")
	err = http.ListenAndServe(":8080", route.ServeHTTP())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
