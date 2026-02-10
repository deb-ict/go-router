package router

import (
	"net/http"
	"testing"
)

func Test_Route_AllowedMethod(t *testing.T) {
	type testCase struct {
		route    *Route
		method   string
		expected int
	}
	tests := []testCase{
		{&Route{methods: nil}, http.MethodGet, 1},
		{&Route{methods: []string{}}, http.MethodGet, 1},
		{&Route{methods: []string{http.MethodGet, http.MethodPost}}, http.MethodGet, 2},
		{&Route{methods: []string{http.MethodGet, http.MethodPost}}, http.MethodDelete, 3},
	}

	for _, tc := range tests {
		result := tc.route.AllowedMethod(tc.method)
		if result != tc.route {
			t.Errorf("Route.AllowedMethod(%s) failed: result not equals instance", tc.method)
		}
		if tc.route.methods == nil {
			t.Errorf("Route.AllowedMethod(%s) failed: method array not initialized", tc.method)
		} else {
			len := len(tc.route.methods)
			if len != tc.expected {
				t.Errorf("Route.AllowedMethod(%s) failed: invalid items in array: got %v, expected %v", tc.method, len, tc.expected)
			}
		}

	}
}

func Test_Route_AllowedMethods(t *testing.T) {
	type testCase struct {
		route    *Route
		methods  []string
		expected int
	}
	tests := []testCase{
		{&Route{methods: nil}, []string{http.MethodGet, http.MethodPost}, 2},
		{&Route{methods: []string{}}, []string{http.MethodGet, http.MethodPost}, 2},
		{&Route{methods: []string{http.MethodGet}}, []string{http.MethodGet, http.MethodPost}, 2},
	}

	for _, tc := range tests {
		result := tc.route.AllowedMethods(tc.methods...)
		if result != tc.route {
			t.Errorf("Route.AllowedMethods(%v) failed: result not equals instance", tc.methods)
		}
		if tc.route.methods == nil {
			t.Errorf("Route.AllowedMethods(%v) failed: method array not initialized", tc.methods)
		} else {
			len := len(tc.route.methods)
			if len != tc.expected {
				t.Errorf("Route.AllowedMethods(%v) failed: invalid items in array: got %v, expected %v", tc.methods, len, tc.expected)
			}
		}
	}
}

func Test_Route_IsMethodAllowed(t *testing.T) {
	type testCase struct {
		route    *Route
		method   string
		expected bool
	}
	tests := []testCase{
		{&Route{methods: nil}, http.MethodGet, true},
		{&Route{methods: []string{}}, http.MethodGet, true},
		{&Route{methods: []string{http.MethodGet, http.MethodPost}}, http.MethodGet, true},
		{&Route{methods: []string{http.MethodGet, http.MethodPost}}, http.MethodDelete, false},
	}

	for _, tc := range tests {
		result := tc.route.IsMethodAllowed(tc.method)
		if result != tc.expected {
			t.Errorf("Route.IsMethodAllowed(%s) failed: got %v, expected %v", tc.method, result, tc.expected)
		}
	}
}
