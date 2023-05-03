package authorization

import (
	"net/http"
)

func WithUnauthorizedHandler(handler http.Handler) MiddlewareOption {
	return func(m *Middleware) {
		m.UnauthorizedHandler = handler
	}
}

func WithUnauthorizedHandlerFunc(handle func(http.ResponseWriter, *http.Request)) MiddlewareOption {
	return func(m *Middleware) {
		m.UnauthorizedHandler = http.HandlerFunc(handle)
	}
}

func WithForbiddenHandler(handler http.Handler) MiddlewareOption {
	return func(m *Middleware) {
		m.ForbiddenHandler = handler
	}
}

func WithForbiddenHandlerFunc(handle func(http.ResponseWriter, *http.Request)) MiddlewareOption {
	return func(m *Middleware) {
		m.ForbiddenHandler = http.HandlerFunc(handle)
	}
}

func WithPolicy(name string, requirements ...Requirement) MiddlewareOption {
	return func(m *Middleware) {
		if m.policies == nil {
			m.policies = make(map[string]Policy)
		}
		m.policies[name] = NewPolicy(name, requirements...)
	}
}
