package router

import (
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

type Middleware interface {
	Middleware(next http.Handler) http.Handler
}

func (mwf MiddlewareFunc) Middleware(next http.Handler) http.Handler {
	return mwf(next)
}

func (r *Router) Use(middlewares ...MiddlewareFunc) {
	if r.middlewares == nil {
		r.middlewares = make([]Middleware, 0)
	}
	for _, m := range middlewares {
		r.middlewares = append(r.middlewares, m)
	}
}

func (r *Router) Middlewares() []Middleware {
	if r.middlewares == nil {
		r.middlewares = make([]Middleware, 0)
	}
	return r.middlewares
}
