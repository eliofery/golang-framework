package middleware

import (
	"github.com/eliofery/golang-image/internal/app/models/user"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userData, err := user.GetCurrentUser(r)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		ctx := user.WithUser(r.Context(), userData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
