package router

import (
	"net/http"
	"testing"
)

type middlewareMock struct {
}

func (m *middlewareMock) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func Test_Router_Use(t *testing.T) {
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}
	router := &Router{
		middlewares: nil,
	}
	router.Use(middleware)

	if router.middlewares == nil {
		t.Error("Router.Use(MiddlewareFunc) failed: slice not initialized")
	}
	if len(router.middlewares) != 1 {
		t.Error("Router.Use(MiddlewareFunc) failed: middleware not appended to slice")
	}
}

func Test_Router_Middleware_ReturnsEmptyList(t *testing.T) {
	router := &Router{
		middlewares: nil,
	}
	result := router.Middlewares()
	if result == nil {
		t.Error("Router.Middlewares() failed: Middleware collection not initialized")
	}
	if len(result) != 0 {
		t.Error("Router.Middlewares() failed: Middleware collection not empty")
	}
}

func Test_Router_Middleware_ReturnsMiddlewares(t *testing.T) {
	middleware := &middlewareMock{}
	router := &Router{
		middlewares: []Middleware{middleware},
	}

	result := router.Middlewares()
	if result == nil {
		t.Error("Router.Middlewares() failed: Middleware collection not empty not initialized")
	}
	if len(result) != 1 {
		t.Error("Router.Middlewares() failed: Middleware not added to collection")
	}
}

func Test_Route_ServerHttp_Middlewares(t *testing.T) {
	executions := make([]int, 0)
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executions = append(executions, 1)
			next.ServeHTTP(w, r)
		})
	}
	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executions = append(executions, 2)
			next.ServeHTTP(w, r)
		})
	}

	router := NewRouter()
	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		executions = append(executions, 3)
	})
	router.Use(middleware1, middleware2)

	req, _ := http.NewRequest(http.MethodGet, "/api", nil)
	rsp := &responseWriterMock{}
	router.ServeHTTP(rsp, req)

	expectedExecutions := 3
	if len(executions) != expectedExecutions {
		t.Errorf("Route.ServerHTTP(middlewares) failed: Incorrect number of executed middlewares: got %v, expected %v", len(executions), expectedExecutions)
	}
	for i := 0; i < expectedExecutions; i++ {
		if len(executions) > i && executions[i] != i+1 {
			t.Errorf("Route.ServerHTTP(middlewares) failed: Incorrect order of execution: got %v, expected %v", executions[i], i+1)
		}
	}
}
