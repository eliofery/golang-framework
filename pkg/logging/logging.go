package logging

import (
	"log/slog"
	"os"
)

var logging *slog.Logger

func New(env string) *slog.Logger {
	var slHandler slog.Handler

	switch env {
	case "dev":
		slHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case "prod":
		slHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logging = slog.New(slHandler)

	return logging
}
