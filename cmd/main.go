package main

import (
	"database/sql"
	"fmt"
	mw "github.com/eliofery/golang-image/internal/app/http/middleware"
	"github.com/eliofery/golang-image/pkg/config"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/database/postgres"
	"github.com/eliofery/golang-image/pkg/email"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/rand"
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

	// Тестовый роут
	route.Get("/", index)
	route.Post("/", post)

	// Запуск сервера
	logger.Info("Сервер запущен: http://localhost:8080")
	err = http.ListenAndServe(":8080", route.ServeHTTP())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func index(ctx router.Ctx) error {
	w := router.ResponseWriter(ctx)
	r := router.Request(ctx)
	l := logging.Logging(ctx)
	db := database.CtxDatabase(ctx)

	var value string
	res := db.QueryRow("SELECT value FROM info WHERE id = 2")
	err := res.Scan(&value)
	if err != nil {
		return err
	}
	fmt.Println(value)

	// Получение токена
	token, err := rand.SessionToken()
	if err != nil {
		l.Info("не удалось получить токен", err)
	}

	data := struct {
		Test string
	}{
		Test: token,
	}

	// Чтение куки
	ck, _ := cookie.Get(r, "test")
	if ck != "" {
		fmt.Println(res)
	} else {
		fmt.Println("нет куки")
	}

	// Удаление куки
	cookie.Delete(w, "test")

	return tpl.Render(ctx, "home", data)
}

func post(ctx router.Ctx) error {
	w := router.ResponseWriter(ctx)
	r := router.Request(ctx)
	l := logging.Logging(ctx)

	// Добавление куки
	cookie.Set(w, "test", "2685723587236582730")

	value := r.FormValue("test")

	// Отправка почты
	emailService, err := email.New()
	if err != nil {
		l.Info("не удалось создать подключение к smtp", err)
	}
	_ = emailService

	mail := email.Email{
		From:    "support@example.kz",
		To:      "guest@example.kz",
		Subject: "Регистрация на сайте",
		Plaintext: `
	       Регистрация прошла успешно.

	       Добро пожаловать к нам на сайт, рады вас видеть.
	   `,
		HTML: `
	       <h1>Регистрация прошла успешно.</h1>

	       <p>Добро пожаловать к нам на сайт, рады вас видеть.</p>
	   `,
	}
	_ = mail
	//err = emailService.Send(mail)
	//if err != nil {
	//   l.Info("не удалось отправить почту", err)
	//}

	_, err = w.Write([]byte(value))
	if err != nil {
		return err
	}

	return nil
}
