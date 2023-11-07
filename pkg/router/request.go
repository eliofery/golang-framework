package router

import (
	"context"
	"net/http"
)

const requestKey key = "request"

func WithRequest(ctx context.Context, w *http.Request) context.Context {
	return context.WithValue(ctx, requestKey, w)
}

func Request(ctx context.Context) *http.Request {
	val := ctx.Value(requestKey)

	request, ok := val.(*http.Request)
	if !ok {
		return nil
	}

	return request
}
