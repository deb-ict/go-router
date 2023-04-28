package authentication

import (
	"net/http"
)

type AuthenticationMiddlewareOption func(*AuthenticationMiddleware)

type AuthenticationMiddleware struct {
	Handler AuthenticationHandler
}

func NewAuthenticationMiddleware(handler AuthenticationHandler, opts ...AuthenticationMiddlewareOption) *AuthenticationMiddleware {
	m := &AuthenticationMiddleware{
		Handler: handler,
	}
	for _, opt := range opts {
		opt(m)
	}
	m.EnsureDefaults()

	return m
}

func (m *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.Handler == nil {
			next.ServeHTTP(w, r)
		} else {
			auth := m.Handler.HandleAuthentication(r)
			ctx := SetAuthenticationContext(r.Context(), auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func (m *AuthenticationMiddleware) EnsureDefaults() {

}
