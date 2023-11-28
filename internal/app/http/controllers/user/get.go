package user

import (
	"github.com/eliofery/golang-image/internal/app/models/user"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
)

func Index(ctx router.Ctx) error {
	userData := user.CtxUser(ctx)

	return tpl.Render(ctx, "user/index", tpl.Data{
		Data: userData,
	})
}

func SignUp(ctx router.Ctx) error {
	return tpl.Render(ctx, "user/signup", tpl.Data{})
}

func SignIn(ctx router.Ctx) error {
	return tpl.Render(ctx, "user/signin", tpl.Data{})
}
