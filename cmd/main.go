package main

import (
	"fmt"
	"github.com/eliofery/golang-image/pkg/config"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/router"
	"log"
	"net/http"
)

func main() {
	// Подключение конфигурационного файла .env
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	// Создание логирования
	log := logging.New("prod")
	_ = log

	// Создание роутера
	route := router.New()

	// Пользовательский Middleware
	route.Use(Middleware)

	// Тестовый роут
	route.Get("/", index)

	// Запуск сервера
	log.Info("Сервер запущен: http://localhost:8080")
	http.ListenAndServe(":8080", route.ServeHTTP())
}

func index(ctx router.Ctx) error {
	w := router.GetResponse(ctx)
	r := router.GetRequest(ctx)
	_ = r

	fmt.Println(w, r)

	w.Write([]byte("Hello World 2"))

	return nil
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("first middleware")

		next.ServeHTTP(w, r)
	})
}
