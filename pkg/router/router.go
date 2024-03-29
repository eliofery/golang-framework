package router

import (
	"context"
	"github.com/eliofery/golang-framework/pkg/logging"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Router struct {
	Mux *chi.Mux
}

type Ctx context.Context
type HandleCtx func(ctx Ctx) error

func New() *Router {
	return &Router{
		Mux: chi.NewRouter(),
	}
}

func (rt *Router) handlerCtx(handler HandleCtx, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if ResponseWriter(ctx) == nil {
		ctx = WithResponseWriter(ctx, w)
	}

	if Request(ctx) == nil {
		ctx = WithRequest(ctx, r)
	}

	r = r.WithContext(ctx)

	if err := handler(r.Context()); err != nil {
		logging.Logging(ctx).Error(err.Error())

		// TODO сделать вывод ошибки через html шаблон
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (rt *Router) Get(path string, handler HandleCtx) {
	rt.Mux.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

func (rt *Router) Post(path string, handler HandleCtx) {
	rt.Mux.Post(path, func(w http.ResponseWriter, r *http.Request) {
		rt.handlerCtx(handler, w, r)
	})
}

func (rt *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	rt.Mux.Use(middlewares...)
}

func (rt *Router) Route(pattern string, fn func(r *Router)) *chi.Mux {
	subRouter := newRouter()

	fn(subRouter)
	rt.Mux.Mount(pattern, subRouter.Mux)

	return subRouter.Mux
}

func newRouter() *Router {
	return &Router{
		Mux: chi.NewRouter(),
	}
}

func (rt *Router) ServeHTTP() http.HandlerFunc {
	return rt.Mux.ServeHTTP
}
