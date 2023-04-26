package router

import (
	"context"
	"net/http"
)

type ContextKey string
type RouteQuery map[string][]string
type RouteParams map[string]string

const (
	queryKey  ContextKey = "router::query"
	paramsKey ContextKey = "router::params"
)

type Router struct {
	tree *Node
}

func NewRouter() *Router {
	r := &Router{
		tree: &Node{},
	}
	return r
}

func Query(r *http.Request) RouteQuery {
	query := r.Context().Value(queryKey)
	if query != nil {
		return query.(RouteQuery)
	}
	return make(RouteQuery)
}

func QueryValues(r *http.Request, key string) []string {
	return Query(r)[key]
}

func QueryValue(r *http.Request, key string) string {
	values := QueryValues(r, key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func Params(r *http.Request) RouteParams {
	params := r.Context().Value(paramsKey)
	if params != nil {
		return params.(RouteParams)
	}
	return make(RouteParams)
}

func Param(r *http.Request, key string) string {
	return Params(r)[key]
}

func (r *Router) HandleFunc(pattern string, handle http.HandlerFunc, opts ...RouteOption) *Route {
	return r.Handle(pattern, http.HandlerFunc(handle), opts...)
}

func (r *Router) Handle(pattern string, handler http.Handler, opts ...RouteOption) *Route {
	node := r.tree.BuildTree(pattern)
	if node == nil {
		return nil
	}
	route := &Route{
		node:    node,
		Handler: handler,
	}
	for _, opt := range opts {
		opt(route)
	}

	if node.Routes == nil {
		node.Routes = make([]*Route, 0)
	}
	node.Routes = append(node.Routes, route)
	return route
}

func (r *Router) findRoute(pattern string, params RouteParams) []*Route {
	node := r.tree.FindNode(pattern, params)
	if node == nil {
		return nil
	}
	if node.Routes == nil || len(node.Routes) == 0 {
		return nil
	}
	return node.Routes
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	params := make(RouteParams)
	routes := r.findRoute(path, params)
	if routes == nil || len(routes) == 0 {
		http.NotFound(w, req)
		return
	}

	hasHandler := false
	for _, route := range routes {
		if route.Handler == nil {
			continue
		}
		hasHandler = true

		if route.IsMethodAllowed(req.Method) {
			r.serverHandle(route.Handler, params, w, req)
			return
		}
	}

	if hasHandler {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	} else {
		http.NotFound(w, req)
	}
}

func (r *Router) serverHandle(h http.Handler, p RouteParams, w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	ctx := req.Context()
	ctx = context.WithValue(ctx, queryKey, query)
	ctx = context.WithValue(ctx, paramsKey, p)
	h.ServeHTTP(w, req.WithContext(ctx))
	return
}
