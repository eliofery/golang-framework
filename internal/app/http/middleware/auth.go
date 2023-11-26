package middleware

import (
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/cookie"
	"net/http"
)

func Auth(next http.Handler) http.Handler {
	uri := "/signin"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := cookie.Get(r, cookie.Session)
		if err != nil {
			http.Redirect(w, r, uri, http.StatusFound)
			return
		}

		service := user.New(r.Context(), user.User{})
		err = service.Auth(token)
		if err != nil {
			http.Redirect(w, r, uri, http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
