package authorization

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func UnauthorizedHandlerMock(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TestUnauthorizedHandler", http.StatusUnauthorized)
}
func ForbiddenHandlerMock(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TestForbiddenHandler", http.StatusForbidden)
}

func Test_WithUnauthorizedHandler(t *testing.T) {
	middleware := &Middleware{}
	handler := http.HandlerFunc(UnauthorizedHandlerMock)
	option := WithUnauthorizedHandler(handler)
	option(middleware)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://testing", nil)
	middleware.UnauthorizedHandler.ServeHTTP(recorder, request)

	responseBody := strings.ReplaceAll(strings.ReplaceAll(recorder.Body.String(), "\r", ""), "\n", "")
	if responseBody != "TestUnauthorizedHandler" {
		t.Errorf("WithUnauthorizedHandler option failed: Invalid body: got %v, expected TestUnauthorizedHandler", responseBody)
	}
}

func Test_WithUnauthorizedHandlerFunc(t *testing.T) {
	middleware := &Middleware{}
	option := WithUnauthorizedHandlerFunc(UnauthorizedHandlerMock)
	option(middleware)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://testing", nil)
	middleware.UnauthorizedHandler.ServeHTTP(recorder, request)

	responseBody := strings.ReplaceAll(strings.ReplaceAll(recorder.Body.String(), "\r", ""), "\n", "")
	if responseBody != "TestUnauthorizedHandler" {
		t.Errorf("WithUnauthorizedHandlerFunc option failed: Invalid body: got %v, expected TestUnauthorizedHandler", responseBody)
	}
}

func Test_WithForbiddenHandler(t *testing.T) {
	middleware := &Middleware{}
	handler := http.HandlerFunc(ForbiddenHandlerMock)
	option := WithForbiddenHandler(handler)
	option(middleware)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://testing", nil)
	middleware.ForbiddenHandler.ServeHTTP(recorder, request)

	responseBody := strings.ReplaceAll(strings.ReplaceAll(recorder.Body.String(), "\r", ""), "\n", "")
	if responseBody != "TestForbiddenHandler" {
		t.Errorf("WithForbiddenHandler option failed: Invalid body: got %v, expected TestForbiddenHandler", responseBody)
	}
}

func Test_WithForbiddenHandlerFunc(t *testing.T) {
	middleware := &Middleware{}
	option := WithForbiddenHandlerFunc(ForbiddenHandlerMock)
	option(middleware)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://testing", nil)
	middleware.ForbiddenHandler.ServeHTTP(recorder, request)

	responseBody := strings.ReplaceAll(strings.ReplaceAll(recorder.Body.String(), "\r", ""), "\n", "")
	if responseBody != "TestForbiddenHandler" {
		t.Errorf("WithForbiddenHandlerFunc option failed: Invalid body: got %v, expected TestForbiddenHandler", responseBody)
	}
}

func Test_WithPolicy(t *testing.T) {
	middleware := &Middleware{
		policies: nil,
	}
	option := WithPolicy("test")
	option(middleware)

	_, ok := middleware.policies["test"]
	if !ok {
		t.Error("WithPolicy option failed: Policy not added")
	}
}
