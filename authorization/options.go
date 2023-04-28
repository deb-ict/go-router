package authorization

import (
	"net/http"
)

func WithUnauthorizedHandler(handler http.Handler) AuthorizationMiddlewareOption {
	return func(m *AuthorizationMiddleware) {
		m.UnauthorizedHandler = handler
	}
}

func WithUnauthorizedHandlerFunc(handle func(http.ResponseWriter, *http.Request)) AuthorizationMiddlewareOption {
	return func(m *AuthorizationMiddleware) {
		m.UnauthorizedHandler = http.HandlerFunc(handle)
	}
}

func WithForbiddenHandler(handler http.Handler) AuthorizationMiddlewareOption {
	return func(m *AuthorizationMiddleware) {
		m.ForbiddenHandler = handler
	}
}

func WithForbiddenHandlerFunc(handle func(http.ResponseWriter, *http.Request)) AuthorizationMiddlewareOption {
	return func(m *AuthorizationMiddleware) {
		m.ForbiddenHandler = http.HandlerFunc(handle)
	}
}
