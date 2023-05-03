package authentication

import (
	"net/http"

	"github.com/deb-ict/go-router"
)

type MiddlewareOption func(*Middleware)

type Middleware struct {
	Handler Handler
}

func NewMiddleware(handler Handler, opts ...MiddlewareOption) *Middleware {
	m := &Middleware{
		Handler: handler,
	}
	for _, opt := range opts {
		opt(m)
	}
	m.EnsureDefaults()

	return m
}

func UseMiddleware(router *router.Router, handler Handler, opts ...MiddlewareOption) {
	m := NewMiddleware(handler, opts...)
	router.Use(m.Middleware)
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var auth Context = nil
		if m.Handler != nil {
			auth = m.Handler.HandleAuthentication(r)
		}
		if auth == nil {
			auth = AnonymouseContext()
		}
		ctx := SetContext(r.Context(), auth)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) EnsureDefaults() {

}
