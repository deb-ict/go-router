package router

import (
	"net/http"
)

type RouteOption func(*Route)

type Route struct {
	node    *Node
	methods []string
	Handler http.Handler
}

func (r *Route) AllowedMethod(method string) {
	if r.methods == nil {
		r.methods = make([]string, 0)
	}
	for _, m := range r.methods {
		if m == method {
			return
		}
	}
	r.methods = append(r.methods, method)
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
