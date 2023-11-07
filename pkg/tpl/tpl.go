package tpl

import (
	"bytes"
	"context"
	"github.com/eliofery/golang-image/internal/resources"
	"github.com/eliofery/golang-image/pkg/router"
	"html/template"
	"io"
	"net/http"
)

type Tpl struct {
	layout string
	page   string
	parts  []string
}

func New(page string) *Tpl {
	return &Tpl{
		layout: getLayout(layoutDefault),
		page:   getPage(page),
		parts:  getParts(),
	}
}

func (t *Tpl) SetLayout(layout string) *Tpl {
	return &Tpl{
		layout: getLayout(layout),
		page:   t.page,
		parts:  t.parts,
	}
}

func (t *Tpl) Render(ctx context.Context, data any) error {
	w := router.Response(ctx)

	tpl, err := template.ParseFS(resources.Views, t.getAllFiles()...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {

	}

	_, err = io.Copy(w, &buf)
	if err != nil {
		return err
	}

	return nil
}

func Render(ctx context.Context, page string, data any) error {
	t := New(page)
	err := t.Render(ctx, data)

	return err
}
