package tpl

import (
	"github.com/eliofery/golang-image/pkg/errors"
	"github.com/eliofery/golang-image/pkg/logging"
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
)

type funcTemplate any

var (
	funcMap = template.FuncMap{
		"csrfInput": csrfInput,
		"errors":    errorsMsg,
	}
)

func csrfInput(r *http.Request, _ Data) funcTemplate {
	return func() template.HTML {
		return csrf.TemplateField(r)
	}
}

func errorsMsg(r *http.Request, data Data) funcTemplate {
	var (
		ErrSomeWrong = errors.New("что то пошло не так")

		errMessage []string
		pubErr     errors.PublicError
	)

	ctx := r.Context()
	l := logging.Logging(ctx)

	for _, err := range data.Errors {
		if errors.As(err, &pubErr) {
			l.Info(pubErr.Public())

			errMessage = append(errMessage, pubErr.Public())
		} else {
			l.Error(pubErr.Error())

			errMessage = append(errMessage, ErrSomeWrong.Error())
		}
	}

	return func() []string {
		return errMessage
	}
}
