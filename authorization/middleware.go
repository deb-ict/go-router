package authorization

import (
	"net/http"

	"github.com/deb-ict/go-router"
	"github.com/deb-ict/go-router/authentication"
)

type AuthorizationMiddlewareOption func(*AuthorizationMiddleware)

type AuthorizationMiddleware struct {
	UnauthorizedHandler http.Handler
	ForbiddenHandler    http.Handler
}

func NewAuthorizationMiddleware(opts ...AuthorizationMiddlewareOption) *AuthorizationMiddleware {
	m := &AuthorizationMiddleware{}
	for _, opt := range opts {
		opt(m)
	}
	m.EnsureDefaults()

	return m
}

func (m *AuthorizationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := router.CurrentRoute(r)
		auth := authentication.GetAuthenticationContext(r.Context())
		if route != nil && route.IsAuthorized() {
			if auth == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			//TODO: Check for authorization policies and return forbidden
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthorizationMiddleware) EnsureDefaults() {
	if m.UnauthorizedHandler == nil {
		m.UnauthorizedHandler = http.HandlerFunc(UnauthorizedHandler)
	}
	if m.ForbiddenHandler == nil {
		m.ForbiddenHandler = http.HandlerFunc(ForbiddenHandler)
	}
}
