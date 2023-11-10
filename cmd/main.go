package main

import (
	"context"
	"fmt"
	"github.com/eliofery/golang-image/pkg/config"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/database/postgres"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/rand"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"github.com/gorilla/csrf"
	"github.com/pressly/goose/v3"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// Подключение конфигурационного файла .env
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	// Создание логирования
	logger := logging.New("prod")

	//// Подключение к БД SQLite
	//driver := sqlite.New()
	//db, err := database.Init(driver)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//// Миграция БД SQLite
	//if err = database.MigrateFS(db, goose.DialectSQLite3); err != nil {
	//	log.Fatal(err)
	//}

	// Подключение к БД Postgres
	driver, err := postgres.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Init(driver)
	if err != nil {
		log.Fatal(err)
	}

	// Миграция БД Postgres
	if err = database.MigrateFS(db, goose.DialectPostgres); err != nil {
		log.Fatal(err)
	}

	//// Отправка почты
	//emailService, err := email.New()
	//if err != nil {
	//    logger.Info("не удалось создать подключение к smtp", err)
	//}
	//
	//mail := email.Email{
	//    From: "support@example.kz",
	//    To: "guest@example.kz",
	//    Subject: "Регистрация на сайте",
	//    Plaintext: `
	//        Регистрация прошла успешно.
	//
	//        Добро пожаловать к нам на сайт, рады вас видеть.
	//    `,
	//    HTML: `
	//        <h1>Регистрация прошла успешно.</h1>
	//
	//        <p>Добро пожаловать к нам на сайт, рады вас видеть.</p>
	//    `,
	//}
	//err = emailService.Send(mail)
	//if err != nil {
	//    logger.Info("не удалось отправить почту", err)
	//}

	// Получение токена
	token, err := rand.SessionToken()
	if err != nil {
		logger.Info("не удалось получить токен", err)
	}
	fmt.Println(token)

	// Создание роутера
	route := router.New()

	// Пользовательский Middleware
	csrfSecure, _ := strconv.ParseBool(os.Getenv("CSRF_SECURE"))
	route.Use(csrf.Protect([]byte(os.Getenv("CSRF_KEY")), csrf.Secure(csrfSecure)))
	route.Use(Middleware)

	// Ресурсы
	tpl.AssetsFsInit(route.Mux)

	// Тестовый роут
	route.Get("/", index)
	route.Post("/", post)

	// Запуск сервера
	logger.Info("Сервер запущен: http://localhost:8080")
	http.ListenAndServe(":8080", route.ServeHTTP())
}

func index(ctx router.Ctx) error {
	val := ctx.Value("test")
	str := val.(string)

	data := struct {
		Test string
	}{
		Test: str,
	}

	err := tpl.Render(ctx, "home", data)

	return err
}

func post(ctx router.Ctx) error {
	w := router.ResponseWriter(ctx)
	r := router.Request(ctx)

	value := r.FormValue("test")

	w.Write([]byte(value))

	return nil
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("first middleware")

		ctx := context.WithValue(r.Context(), "test", "Тест передачи контекста через middleware прошел успешно")
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
