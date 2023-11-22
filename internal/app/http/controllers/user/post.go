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
		l.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)

		return tpl.Render(ctx, "user/signup", tpl.Data{
			Errors: []error{err},
		})
	}

	http.Redirect(w, r, "/user", http.StatusFound)

	return nil
}
