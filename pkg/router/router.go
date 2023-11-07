package router

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Router struct {
	Mux *chi.Mux
}

type Ctx context.Context
type HandleCtx func(ctx Ctx) error
type key string

func New() *Router {
	return &Router{
		Mux: chi.NewRouter(),
	}
}

func (chi *Router) handlerCtx(handler HandleCtx, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = WithResponse(ctx, w)
	ctx = WithRequest(ctx, r)

	r = r.WithContext(ctx)

	if err := handler(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (chi *Router) Get(path string, handler HandleCtx) {
	chi.Mux.Get(path, func(w http.ResponseWriter, r *http.Request) {
		chi.handlerCtx(handler, w, r)
	})
}

func (chi *Router) Post(path string, handler HandleCtx) {
	chi.Mux.Post(path, func(w http.ResponseWriter, r *http.Request) {
		chi.handlerCtx(handler, w, r)
	})
}

func (chi *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	chi.Mux.Use(middlewares...)
}

func (chi *Router) ServeHTTP() http.HandlerFunc {
	return chi.Mux.ServeHTTP
}
