package user

import (
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"net/http"
)

func Create(ctx router.Ctx) error {
	w, r, l := router.ResponseWriter(ctx), router.Request(ctx), logging.Logging(ctx)

	service := user.New(ctx, user.User{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})

	err := service.SignUp()
	if err != nil {
		l.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/signup", tpl.Data{
			Data:   service.User,
			Errors: []error{err},
		})
	}

	http.Redirect(w, r, "/user", http.StatusFound)

	return nil
}

func Auth(ctx router.Ctx) error {
	w, r, l := router.ResponseWriter(ctx), router.Request(ctx), logging.Logging(ctx)

	service := user.New(ctx, user.User{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})

	err := service.SignIn()
	if err != nil {
		l.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/signin", tpl.Data{
			Data:   service.User,
			Errors: []error{err},
		})
	}

	http.Redirect(w, r, "/user", http.StatusFound)

	return nil
}
