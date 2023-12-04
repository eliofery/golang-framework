package user

import (
	"github.com/eliofery/golang-framework/pkg/router"
	"github.com/eliofery/golang-framework/pkg/tpl"
)

func Index(ctx router.Ctx) error {
	return tpl.Render(ctx, "user/index", tpl.Data{})
}

func SignUp(ctx router.Ctx) error {
	return tpl.Render(ctx, "user/signup", tpl.Data{})
}

func SignIn(ctx router.Ctx) error {
	return tpl.Render(ctx, "user/signin", tpl.Data{})
}

func ForgotPassword(ctx router.Ctx) error {
	return tpl.Render(ctx, "user/forgot-pw", tpl.Data{})
}

func ResetPassword(ctx router.Ctx) error {
	r := router.Request(ctx)

	return tpl.Render(ctx, "user/reset-pw", tpl.Data{
		Data: r.FormValue("token"),
	})
}
