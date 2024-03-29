package main

import (
	"database/sql"
	"github.com/eliofery/golang-framework/internal/app/http/controllers/home"
	"github.com/eliofery/golang-framework/internal/app/http/controllers/user"
	mw "github.com/eliofery/golang-framework/internal/app/http/middleware"
	"github.com/eliofery/golang-framework/pkg/config"
	"github.com/eliofery/golang-framework/pkg/database"
	"github.com/eliofery/golang-framework/pkg/database/postgres"
	"github.com/eliofery/golang-framework/pkg/logging"
	"github.com/eliofery/golang-framework/pkg/router"
	"github.com/eliofery/golang-framework/pkg/tpl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
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

	// Подключение валидатора
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Создание роутера
	route := router.New()

	// Пользовательский Middleware
	route.Use(middleware.RequestID)
	route.Use(middleware.RealIP)
	route.Use(middleware.Logger)
	route.Use(middleware.Recoverer)
	route.Use(middleware.URLFormat)
	route.Use(mw.Csrf)
	route.Use(mw.Inject(logger, db, validate))
	route.Use(mw.SetUser)

	// Подключение ресурсов
	tpl.AssetsFsInit(route.Mux)

	// Роуты
	route.Get("/", home.Index)
	route.Post("/", home.Post)

	route.Route("/user", func(r *router.Router) {
		r.Mux.Use(mw.Auth)

		r.Get("/", user.Index)
		r.Post("/logout", user.Logout)
	})

	route.Route("/", func(r *router.Router) {
		r.Mux.Use(mw.Guest)

		r.Get("/signup", user.SignUp)
		r.Post("/signup", user.Create)

		r.Get("/signin", user.SignIn)
		r.Post("/signin", user.Auth)

		r.Get("/forgot-pw", user.ForgotPassword)
		r.Post("/forgot-pw", user.ProcessForgotPassword)

		r.Get("/reset-pw", user.ResetPassword)
		r.Post("/reset-pw", user.ProcessResetPassword)
	})

	// Запуск сервера
	logger.Info("Сервер запущен: http://localhost:8080")
	err = http.ListenAndServe(":8080", route.ServeHTTP())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
