package tpl

import (
	"github.com/eliofery/golang-image/internal/resources"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	assetsDir = "internal/resources/assets"
	pattern   = "/assets/*"
	prefix    = "/assets/"
)

func AssetsInit(route *chi.Mux) {
	fs := http.FileServer(http.Dir(assetsDir))
	route.Handle(pattern, http.StripPrefix(prefix, fs))
}

func AssetsFsInit(route *chi.Mux) {
	route.Handle(pattern, http.FileServer(http.FS(resources.Assets)))
}
