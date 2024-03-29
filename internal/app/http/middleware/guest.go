package middleware

import (
	"github.com/eliofery/golang-framework/internal/app/models/user"
	"net/http"
)

func Guest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user.CtxUser(r.Context()) != nil {
			http.Redirect(w, r, "/user", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
