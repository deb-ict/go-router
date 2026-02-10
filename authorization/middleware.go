package authorization

import (
	"net/http"

	"github.com/deb-ict/go-router"
	"github.com/deb-ict/go-router/authentication"
)

type MiddlewareOption func(*Middleware)

type Middleware struct {
	UnauthorizedHandler http.Handler
	ForbiddenHandler    http.Handler
	policies            map[string]Policy
}

func NewMiddleware(opts ...MiddlewareOption) *Middleware {
	m := &Middleware{
		policies: make(map[string]Policy),
	}
	for _, opt := range opts {
		opt(m)
	}
	m.EnsureDefaults()

	return m
}

func UseMiddleware(router *router.Router, opts ...MiddlewareOption) {
	m := NewMiddleware(opts...)
	router.Use(m.Middleware)
}

func (m *Middleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := router.CurrentRoute(r)
		if route != nil && route.IsAuthorized() {
			auth := authentication.GetContext(r.Context())
			if !auth.IsAuthenticated() {
				m.UnauthorizedHandler.ServeHTTP(w, r)
				return
			}

			policyName := route.GetAuthorizationPolicy()
			policy, ok := m.policies[policyName]
			if !ok {
				m.UnauthorizedHandler.ServeHTTP(w, r)
				return
			}

			if !policy.MeetsRequirements(auth) {
				m.ForbiddenHandler.ServeHTTP(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) EnsureDefaults() {
	if m.UnauthorizedHandler == nil {
		m.UnauthorizedHandler = http.HandlerFunc(UnauthorizedHandler)
	}
	if m.ForbiddenHandler == nil {
		m.ForbiddenHandler = http.HandlerFunc(ForbiddenHandler)
	}
}

func (m *Middleware) SetPolicy(policy Policy) {
	m.policies[policy.GetName()] = policy
}

func (m *Middleware) GetPolicy(name string) (Policy, bool) {
	policy, ok := m.policies[name]
	return policy, ok
}
