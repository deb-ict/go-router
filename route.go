package router

import (
	"net/http"
)

type RouteOption func(*Route)

type Route struct {
	node           *Node
	methods        []string
	Handler        http.Handler
	authPolicyName string
}

func (r *Route) AllowedMethod(method string) *Route {
	if r.methods == nil {
		r.methods = make([]string, 0)
	}
	for _, m := range r.methods {
		if m == method {
			return r
		}
	}
	r.methods = append(r.methods, method)
	return r
}

func (r *Route) IsMethodAllowed(method string) bool {
	if r.methods == nil || len(r.methods) == 0 {
		return true
	}
	for _, m := range r.methods {
		if m == method {
			return true
		}
	}
	return false
}

func (r *Route) Authorize(policyName string) {
	r.authPolicyName = policyName
}

func (r *Route) IsAuthorized() bool {
	return r.authPolicyName != ""
}

func (r *Route) GetAuthorizationPolicy() string {
	return r.authPolicyName
}
