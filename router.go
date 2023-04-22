package router

import (
	"context"
	"net/http"
)

type contextKey string
type RouteQuery map[string][]string
type RouteParams map[string]string

const (
	queryKey  contextKey = "query"
	paramsKey contextKey = "params"
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
	node.Route = route
	return node.Route
}

func (r *Router) findRoute(pattern string, params RouteParams) *Route {
	node := r.tree.FindNode(pattern, params)
	if node == nil {
		return nil
	}
	return node.Route
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	params := make(RouteParams)
	route := r.findRoute(path, params)
	if route == nil || route.Handler == nil {
		http.NotFound(w, req)
		return
	}
	if !route.IsMethodAllowed(req.Method) {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := req.URL.Query()

	ctx := req.Context()
	ctx = context.WithValue(ctx, queryKey, query)
	ctx = context.WithValue(ctx, paramsKey, params)
	route.Handler.ServeHTTP(w, req.WithContext(ctx))
}
