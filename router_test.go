package router

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

type responseWriterMock struct {
	statusCode int
}

func (mock *responseWriterMock) Header() http.Header {
	return http.Header{}
}

func (mock *responseWriterMock) Write([]byte) (int, error) {
	return 0, nil
}

func (mock *responseWriterMock) WriteHeader(statusCode int) {
	mock.statusCode = statusCode
}

func Test_NewRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Errorf("NewRouter() failed: instance is <nil>")
	}
	if router != nil && router.tree == nil {
		t.Errorf("NewRouter() failed: instance.tree is <nil>")
	}
}

func Test_Query_Default(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	query := Query(req)
	if query == nil {
		t.Error("Query() failed: got <nil>, expected <empty>")
	}
}

func Test_Params_Default(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	params := Params(req)
	if params == nil {
		t.Error("Params() failed: got <nil>, expected <empty>")
	}
}

func Test_CurrentRoute(t *testing.T) {
	type testCase struct {
		route *Route
	}
	tests := []testCase{
		{&Route{}},
		{nil},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		if tc.route != nil {
			ctx := context.WithValue(req.Context(), routeKey, tc.route)
			req = req.WithContext(ctx)
		}

		result := CurrentRoute(req)
		if tc.route == nil && result != nil {
			t.Error("CurrentRoute(nil) failed: result is not nil")
		} else if tc.route != nil && result == nil {
			t.Error("CurrentRoute(route) failed: result is nil")
		} else if tc.route != nil && result != tc.route {
			t.Error("CurrentRoute(route) failed: result not equals route")
		}
	}

}

func Test_QueryValues(t *testing.T) {
	type testCase struct {
		lookup   string
		expected []string
	}
	tests := []testCase{
		{"test", []string{"ok"}},
		{"notfound", []string{}},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		query := make(RouteQuery)
		query["test"] = []string{"ok"}
		ctx := context.WithValue(req.Context(), queryKey, query)
		req = req.WithContext(ctx)

		result := QueryValues(req, tc.lookup)
		if len(result) != len(tc.expected) {
			t.Errorf("QueryValues(%s) failed: got %v items, expected %v items", tc.lookup, len(result), len(tc.expected))
		}
		if len(tc.expected) > 0 && result[0] != tc.expected[0] {
			t.Errorf("QueryValues(%s) failed: got %v, expected %v", tc.lookup, result[0], tc.expected[0])
		}
	}
}

func Test_QueryValue(t *testing.T) {
	type testCase struct {
		lookup   string
		expected string
	}
	tests := []testCase{
		{"test", "ok"},
		{"notfound", ""},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		query := make(RouteQuery)
		query["test"] = []string{"ok", "second"}
		ctx := context.WithValue(req.Context(), queryKey, query)
		req = req.WithContext(ctx)

		result := QueryValue(req, tc.lookup)
		if result != tc.expected {
			t.Errorf("QueryValue(%s) failed: got %v, expected %v", tc.lookup, result[0], tc.expected[0])
		}
	}
}

func Test_Param(t *testing.T) {
	type testCase struct {
		lookup   string
		expected string
	}
	tests := []testCase{
		{"test", "ok"},
		{"notfound", ""},
	}

	for _, tc := range tests {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		params := make(RouteParams)
		params["test"] = "ok"
		ctx := context.WithValue(req.Context(), paramsKey, params)
		req = req.WithContext(ctx)

		result := Param(req, tc.lookup)
		if result != tc.expected {
			t.Errorf("Param(%s) failed: got %v, expected %v", tc.lookup, result, tc.expected)
		}
	}
}

func Test_Router_HandleFunc(t *testing.T) {
	type testCase struct {
		pattern     string
		expectedNil bool
	}
	tests := []testCase{
		{"", true},
		{"api", true},
		{"/api", false},
	}

	for _, tc := range tests {
		optionApplied := false
		option := func(r *Route) {
			optionApplied = true
		}
		handle := func(http.ResponseWriter, *http.Request) {}
		router := &Router{tree: &Node{}}
		result := router.HandleFunc(tc.pattern, handle, option)
		if tc.expectedNil {
			if result != nil {
				t.Errorf("Router.HandleFunc(%s) failed: expected <nil>", tc.pattern)
			}
		} else {
			if result == nil {
				t.Errorf("Router.HandleFunc(%s) failed: expected instance", tc.pattern)
			} else {
				if result.node == nil {
					t.Errorf("Router.HandleFunc(%s) failed: expected instance.node not <nil>", tc.pattern)
				}
				if result.handler == nil {
					t.Errorf("Router.HandleFunc(%s) failed: expected instance.handler not <nil>", tc.pattern)
				}
				if !optionApplied {
					t.Errorf("Router.HandleFunc(%s) failed: option not applied", tc.pattern)
				}
			}
		}
	}
}

func Test_Router_Handle(t *testing.T) {
	type testCase struct {
		pattern     string
		expectedNil bool
	}
	tests := []testCase{
		{"", true},
		{"api", true},
		{"/api", false},
	}

	for _, tc := range tests {
		optionApplied := false
		option := func(r *Route) {
			optionApplied = true
		}
		handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
		router := &Router{tree: &Node{}}
		result := router.Handle(tc.pattern, handler, option)
		if tc.expectedNil {
			if result != nil {
				t.Errorf("Router.Handle(%s) failed: expected <nil>", tc.pattern)
			}
		} else {
			if result == nil {
				t.Errorf("Router.Handle(%s) failed: expected instance", tc.pattern)
			} else {
				if result.node == nil {
					t.Errorf("Router.Handle(%s) failed: expected instance.node not <nil>", tc.pattern)
				}
				if result.handler == nil {
					t.Errorf("Router.Handle(%s) failed: expected instance.handler not <nil>", tc.pattern)
				}
				if !optionApplied {
					t.Errorf("Router.Handle(%s) failed: option not applied", tc.pattern)
				}
			}
		}
	}
}

func Test_Router_findRoute(t *testing.T) {
	type testCase struct {
		pattern     string
		expectedNil bool
	}
	tests := []testCase{
		{"/api", false},
		{"/nilroutes", true},
		{"/emptyroutes", true},
		{"/notfound", true},
	}

	for _, tc := range tests {
		router := &Router{
			tree: &Node{
				Nodes: []*Node{
					{
						Segment: "api",
						Type:    NodeTypePath,
						Routes:  []*Route{{}},
					},
					{
						Segment: "nilroutes",
						Type:    NodeTypePath,
						Routes:  nil,
					},
					{
						Segment: "emptyroutes",
						Type:    NodeTypePath,
						Routes:  []*Route{},
					},
				},
			},
		}

		params := make(RouteParams)
		result := router.findRoute(tc.pattern, params)
		if tc.expectedNil && result != nil {
			t.Errorf("Router.findNode(%s) failed: Expected <nil>", tc.pattern)
		}
		if !tc.expectedNil && result == nil {
			t.Errorf("Router.findNode(%s) failed: Expected not <nil>", tc.pattern)
		}
	}
}

func Test_Router_ServeHttp(t *testing.T) {
	type testCase struct {
		pattern  string
		expected int
	}
	handle := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	tests := []testCase{
		{"/api", 200},
		{"/test", 404},
		{"/get", 405},
		{"/nohandler", 404},
	}

	router := NewRouter()
	router.HandleFunc("/api", handle)
	router.HandleFunc("/get", handle).AllowedMethod(http.MethodPost)
	router.Handle("/nohandler", nil)
	for _, tc := range tests {
		req, _ := http.NewRequest(http.MethodGet, tc.pattern, nil)
		rsp := &responseWriterMock{}
		router.ServeHTTP(rsp, req)

		if rsp.statusCode != tc.expected {
			t.Errorf("Router.ServeHTTP(%s) failed: invalid status code: got %v, expected %v", tc.pattern, rsp.statusCode, tc.expected)
		}
	}
}

func Test_Router_ServeHttp_SubRouter(t *testing.T) {
	middlewareCalls := make([]string, 0)
	handle := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	rootMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalls = append(middlewareCalls, "root")
			h.ServeHTTP(w, r)
		})
	}
	subMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			middlewareCalls = append(middlewareCalls, "sub")
			h.ServeHTTP(w, r)
		})
	}

	router := NewRouter()
	router.Use(rootMiddleware)

	subRouter := router.PathPrefix("/api/v1").SubRouter()
	subRouter.Use(subMiddleware)
	subRouter.HandleFunc("/items", handle)

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/items", nil)
	rsp := &responseWriterMock{}
	router.ServeHTTP(rsp, req)

	if rsp.statusCode != http.StatusOK {
		t.Errorf("Router.ServeHTTP(%s) failed: invalid status code: got %v, expected %v", "/api/v1/items", rsp.statusCode, http.StatusOK)
	}
	expectedMiddlewareCalls := []string{"root", "sub"}
	if !reflect.DeepEqual(middlewareCalls, expectedMiddlewareCalls) {
		t.Errorf("Router.ServeHTTP(%s) failed: middleware calls = %v, expected %v", "/api/v1/items", middlewareCalls, expectedMiddlewareCalls)
	}
}

func Test_Router_ServeHttp_MultiMethod_SingleHandler(t *testing.T) {
	type testCase struct {
		pattern  string
		method   string
		expected int
	}
	tests := []testCase{
		{"/api", http.MethodGet, 200},
		{"/api", http.MethodPut, 200},
		{"/api", http.MethodDelete, 405},
	}
	handlerCalled := 0
	handlerCalledExpected := 2
	handle := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled++
		w.WriteHeader(http.StatusOK)
	}

	router := NewRouter()
	router.HandleFunc("/api", handle).AllowedMethod(http.MethodGet)
	router.HandleFunc("/api", handle).AllowedMethod(http.MethodPut)
	for _, tc := range tests {
		req, _ := http.NewRequest(tc.method, tc.pattern, nil)
		rsp := &responseWriterMock{}
		router.ServeHTTP(rsp, req)

		if rsp.statusCode != tc.expected {
			t.Errorf("Router.ServeHTTP(%s) failed: invalid status code: got %v, expected %v", tc.method, rsp.statusCode, tc.expected)
		}
	}

	if handlerCalled != handlerCalledExpected {
		t.Errorf("Router.ServerHTTP() failed: handler called %v times, expected %v", handlerCalled, handlerCalledExpected)
	}
}

func Test_Router_ServeHttp_MultiMethod_MultiHandler(t *testing.T) {
	type testCase struct {
		pattern  string
		method   string
		expected int
	}
	tests := []testCase{
		{"/api", http.MethodGet, 200},
		{"/api", http.MethodPut, 200},
		{"/api", http.MethodDelete, 405},
	}

	getHandlerCalled := 0
	getHandle := func(w http.ResponseWriter, r *http.Request) {
		getHandlerCalled++
		w.WriteHeader(http.StatusOK)
	}
	putHandlerCalled := 0
	putHandle := func(w http.ResponseWriter, r *http.Request) {
		putHandlerCalled++
		w.WriteHeader(http.StatusOK)
	}

	router := NewRouter()
	router.HandleFunc("/api", getHandle).AllowedMethod(http.MethodGet)
	router.HandleFunc("/api", putHandle).AllowedMethod(http.MethodPut)
	for _, tc := range tests {
		req, _ := http.NewRequest(tc.method, tc.pattern, nil)
		rsp := &responseWriterMock{}
		router.ServeHTTP(rsp, req)

		if rsp.statusCode != tc.expected {
			t.Errorf("Router.ServeHTTP(%s) failed: invalid status code: got %v, expected %v", tc.method, rsp.statusCode, tc.expected)
		}
	}

	if getHandlerCalled != 1 {
		t.Errorf("Router.ServerHTTP(%s) failed: handler called %v times, expected 1", http.MethodGet, getHandlerCalled)
	}
	if putHandlerCalled != 1 {
		t.Errorf("Router.ServerHTTP(%s) failed: handler called %v times, expected 1", http.MethodPut, putHandlerCalled)
	}
}
