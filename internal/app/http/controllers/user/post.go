package user

import (
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
	"net/http"
)

func Create(ctx router.Ctx) error {
	w := router.ResponseWriter(ctx)
	r := router.Request(ctx)
	l := logging.Logging(ctx)

	modelUser := user.New(ctx)
	err := modelUser.Create()
	if err != nil {
		l.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/signup", tpl.Data{
			Data:   modelUser.Dto,
			Errors: []error{err},
		})
	}

	http.Redirect(w, r, "/user", http.StatusFound)

	return nil
}

func Auth(ctx router.Ctx) error {
	w := router.ResponseWriter(ctx)
	r := router.Request(ctx)
	l := logging.Logging(ctx)

	modelUser := user.New(ctx)
	err := modelUser.Auth()
	if err != nil {
		l.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/signin", tpl.Data{
			Data:   modelUser.Dto,
			Errors: []error{err},
		})
	}

	http.Redirect(w, r, "/user", http.StatusFound)

	return nil
}
