package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var (
	ErrNotFound = fmt.Errorf("не удалось загрузить файл .env")
)

const (
	Dev = ".env"
)

func Init() error {
	var env string

	if len(os.Args) == 1 {
		env = Dev
	} else {
		env = os.Args[1]
	}

	if err := godotenv.Load(env); err != nil {
		return ErrNotFound
	}

	return nil
}
