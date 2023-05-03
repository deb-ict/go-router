package validation

import (
	"net/http"
)

type MiddlewareOption func(*Middleware)

type Middleware struct {
}

func NewMiddleware(opts ...MiddlewareOption) *Middleware {
	m := &Middleware{}
	for _, opt := range opts {
		opt(m)
	}
	m.EnsureDefaults()

	return m
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var validation Context

		ctx := SetContext(r.Context(), validation)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) EnsureDefaults() {

}
