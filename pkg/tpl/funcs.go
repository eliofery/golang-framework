package tpl

import (
	"github.com/gorilla/csrf"
	"html/template"
	"net/http"
)

type funcTemplate func() template.HTML

var (
	funcMap = template.FuncMap{
		"csrfInput": csrfInput,
	}
)

func csrfInput(r *http.Request) funcTemplate {
	return func() template.HTML {
		return csrf.TemplateField(r)
	}
}
