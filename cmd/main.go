package main

import (
	"github.com/eliofery/golang-image/pkg/config"
	"github.com/eliofery/golang-image/pkg/logging"
	"log"
	"os"
)

func main() {
	// Подключение конфигурационного файла .env
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	// Создание логирования
	log := logging.New("prod")

	// Тестирование лога и переменного окружения
	log.Info(os.Getenv("ENV"))
}
