package tpl

import (
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
)

type funcTemplate any

var (
	funcMap = template.FuncMap{
		"csrfInput": csrfInput,
		"errors":    errors,
	}
)

func csrfInput(r *http.Request, _ Data) funcTemplate {
	return func() template.HTML {
		return csrf.TemplateField(r)
	}
}

func errors(_ *http.Request, data Data) funcTemplate {
	return func() []string {
		return data.Errors
	}
}
