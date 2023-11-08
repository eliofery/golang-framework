package tpl

import (
	"bytes"
	"context"
	"github.com/eliofery/golang-image/internal/resources"
	"github.com/eliofery/golang-image/pkg/router"
	"html/template"
	"io"
	"path"
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
	w := router.ResponseWriter(ctx)
	r := router.Request(ctx)

	layoutFileName := path.Base(t.getAllFiles()[0])
	tpl := template.New(layoutFileName)
	tpl = tpl.Funcs(getFuncMap(r))

	tpl, err := tpl.ParseFS(resources.Views, t.getAllFiles()...)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, data); err != nil {
		return err
	}

	if _, err = io.Copy(w, &buf); err != nil {
		return err
	}

	return nil
}

func Render(ctx context.Context, page string, data any) error {
	t := New(page)
	err := t.Render(ctx, data)

	return err
}
