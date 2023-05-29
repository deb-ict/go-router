package authorization

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_UnauthorizedHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://testing", nil)
	handler := http.HandlerFunc(UnauthorizedHandler)
	handler.ServeHTTP(recorder, request)

	responseCode := recorder.Code
	responseBody := strings.ReplaceAll(strings.ReplaceAll(recorder.Body.String(), "\r", ""), "\n", "")
	if responseCode != http.StatusUnauthorized {
		t.Errorf("UnauthorizedHandler failed: Invalid status code: got %v, expected %v", responseCode, http.StatusUnauthorized)
	}
	if responseBody != "Unauthorized" {
		t.Errorf("UnauthorizedHandler failed: Invalid body: got %v, expected Unauthorized", responseBody)
	}
}

func Test_ForbiddenHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "http://testing", nil)
	handler := http.HandlerFunc(ForbiddenHandler)
	handler.ServeHTTP(recorder, request)

	responseCode := recorder.Code
	responseBody := strings.ReplaceAll(strings.ReplaceAll(recorder.Body.String(), "\r", ""), "\n", "")
	if responseCode != http.StatusForbidden {
		t.Errorf("ForbiddenHandler failed: Invalid status code: got %v, expected %v", responseCode, http.StatusForbidden)
	}
	if responseBody != "Forbidden" {
		t.Errorf("ForbiddenHandler failed: Invalid body: got %v, expected Forbidden", responseBody)
	}
}
