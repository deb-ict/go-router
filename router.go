package router

import (
	"context"
	"net/http"
	"strings"
)

type ContextKey string
type RouteQuery map[string][]string
type RouteParams map[string]string

const (
	routeKey  ContextKey = "router:route"
	queryKey  ContextKey = "router::query"
	paramsKey ContextKey = "router::params"
)

type Router struct {
	tree        *Node
	middlewares []Middleware
}

func NewRouter() *Router {
	r := &Router{
		tree:        &Node{},
		middlewares: make([]Middleware, 0),
	}
	return r
}

func CurrentRoute(r *http.Request) *Route {
	value := r.Context().Value(routeKey)
	if value == nil {
		return nil
	}
	return value.(*Route)
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
	normalizedKey := strings.ToLower(key)
	return Params(r)[normalizedKey]
}

func (r *Router) PathPrefix(pattern string, opts ...RouteOption) *Route {
	node := r.tree.BuildTree(pattern)
	if node == nil {
		return nil
	}
	route := &Route{
		node: node,
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

func (r *Router) HandleFunc(pattern string, handle http.HandlerFunc, opts ...RouteOption) *Route {
	return r.Handle(pattern, http.HandlerFunc(handle), opts...)
}

func (r *Router) Handle(pattern string, handler http.Handler, opts ...RouteOption) *Route {
	route := r.PathPrefix(pattern, opts...)
	if route != nil {
		route.Handler = handler
	}
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
			r.serverRoute(route, params, w, req)
			return
		}
	}

	if hasHandler {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	} else {
		http.NotFound(w, req)
	}
}

func (r *Router) serverRoute(route *Route, params RouteParams, w http.ResponseWriter, req *http.Request) {
	urlValues := req.URL.Query()
	query := make(RouteQuery)
	for key, values := range urlValues {
		query[key] = values
	}

	ctx := req.Context()
	ctx = context.WithValue(ctx, routeKey, route)
	ctx = context.WithValue(ctx, queryKey, query)
	ctx = context.WithValue(ctx, paramsKey, params)

	handler := route.Handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i].Middleware(handler)
	}
	handler.ServeHTTP(w, req.WithContext(ctx))
}
