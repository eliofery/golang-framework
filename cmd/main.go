package main

import (
	"context"
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
	w := router.Response(ctx)
	r := router.GetRequest(ctx)

	val := ctx.Value("test")
	str, ok := val.(string)
	if !ok {
		return nil
	}

	_ = r

	w.Write([]byte("Hello World " + str))

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
