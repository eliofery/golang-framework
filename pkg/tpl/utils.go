package tpl

import (
	"github.com/eliofery/golang-image/internal/resources"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
)

const (
	layoutDefault = "default"
	fileExt       = ".html"
)

var (
	pathView    = pathJoin("views")
	pathLayouts = pathJoin(pathView + "/layouts")
	pathPages   = pathJoin(pathView + "/pages")
	pathParts   = pathJoin(pathView + "/parts")
)

func pathJoin(p string) string {
	return path.Join(strings.Split(p, "/")...)
}

func getLayout(layout string) string {
	return path.Join(pathLayouts, layout+fileExt)
}

func getPage(page string) string {
	return path.Join(pathPages, pathJoin(page)+fileExt)
}

func getParts() []string {
	parts, err := fs.ReadDir(resources.Views, pathParts)
	if err != nil {
		log.Println("Error reading folder parts:", err)
	}

	partsNew := make([]string, 0, len(parts))
	for _, part := range parts {
		partsNew = append(partsNew, path.Join(pathParts, part.Name()))
	}

	return partsNew
}

func (t *Tpl) getAllFiles() []string {
	files := []string{t.layout, t.page}
	files = append(files, t.parts...)

	return files
}

func getFuncMap(r *http.Request) template.FuncMap {
	var fMap = template.FuncMap{}

	for key, callback := range funcMap {
		cb, ok := callback.(func(*http.Request) funcTemplate)
		if !ok {
			continue
		}

		fMap[key] = cb(r)
	}

	return fMap
}
