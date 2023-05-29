package authentication

import (
	"testing"
)

func Test_WithApiKeyAuthenticationHeaderName(t *testing.T) {
	expected := "test"
	handler := &ApiKeyAuthenticationHandler{}
	option := WithApiKeyAuthenticationHeaderName(expected)
	option(handler)

	if handler.HeaderName != expected {
		t.Errorf("WithApiKeyAuthenticationHeaderName() failed: got %s, expected %s", handler.HeaderName, expected)
	}
}

func Test_WithApiKeyAuthenticationQueryParamName(t *testing.T) {
	expected := "test"
	handler := &ApiKeyAuthenticationHandler{}
	option := WithApiKeyAuthenticationQueryParamName(expected)
	option(handler)

	if handler.QueryParamName != expected {
		t.Errorf("WithApiKeyAuthenticationQueryParamName() failed: got %s, expected %s", handler.QueryParamName, expected)
	}
}
