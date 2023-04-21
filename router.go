package router

import (
	"net/http"
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

func (r *Router) HandleFunc(pattern string, handle http.HandlerFunc, opts ...RouteOption) *Route {
	return r.Handle(pattern, http.HandlerFunc(handle), opts...)
}

func (r *Router) Handle(pattern string, handler http.Handler, opts ...RouteOption) *Route {
	route := &Route{
		Handler: handler,
	}
	for _, opt := range opts {
		opt(route)
	}

	node := r.tree.BuildTree(pattern)
	node.Route = route
	return node.Route
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	node := r.tree.FindNode(path)
	if node == nil || node.Route == nil || node.Route.Handler == nil {
		http.NotFound(w, req)
		return
	}

	node.Route.Handler.ServeHTTP(w, req)
}
