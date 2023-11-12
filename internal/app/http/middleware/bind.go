package middleware

import (
	"database/sql"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/logging"
	"log/slog"
	"net/http"
)

func Bind(logger *slog.Logger, db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logging.WithLogging(r.Context(), logger)
			ctx = database.WithDatabase(ctx, db)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
