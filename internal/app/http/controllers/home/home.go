package home

import (
	"fmt"
	"github.com/eliofery/golang-image/pkg/cookie"
	"github.com/eliofery/golang-image/pkg/database"
	"github.com/eliofery/golang-image/pkg/email"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/eliofery/golang-image/pkg/rand"
	"github.com/eliofery/golang-image/pkg/router"
	"github.com/eliofery/golang-image/pkg/tpl"
)

func Index(ctx router.Ctx) error {
	w := router.ResponseWriter(ctx)
	r := router.Request(ctx)
	l := logging.Logging(ctx)
	db := database.CtxDatabase(ctx)
	_ = db

	// Получение токена
	token, err := rand.SessionToken()
	if err != nil {
		l.Info("не удалось получить токен", err)
	}

	data := tpl.Data{
		Meta: tpl.Meta{
			Title: "Главная",
		},
		Data: struct {
			Token string
		}{
			Token: token,
		},
		Errors: tpl.PublicErrors(
			"ошибка 1",
			"ошибка 2",
			"ошибка 3",
		),
	}

	// Чтение куки
	ck, _ := cookie.Get(r, "test")
	fmt.Println(ck)

	// Удаление куки
	cookie.Delete(w, "test")

	return tpl.Render(ctx, "home", data)
}

func Post(ctx router.Ctx) error {
	op := "postHandler"

	w := router.ResponseWriter(ctx)
	r := router.Request(ctx)

	// Добавление куки
	cookie.Set(w, "test", "2685723587236582730")

	value := r.FormValue("test")

	// Отправка почты
	emailService := email.New()
	_ = emailService

	mail := email.Email{
		From:    "support@example.kz",
		To:      "guest@example.kz",
		Subject: "Регистрация на сайте",
		Plaintext: `
	       Регистрация прошла успешно.

	       Добро пожаловать к нам на сайт, рады вас видеть.
	   `,
		HTML: `
	       <h1>Регистрация прошла успешно.</h1>

	       <p>Добро пожаловать к нам на сайт, рады вас видеть.</p>
	   `,
	}
	_ = mail
	//err = emailService.Send(mail)
	//if err != nil {
	//   l.Info("не удалось отправить почту", err)
	//}

	_, err := w.Write([]byte(value))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
