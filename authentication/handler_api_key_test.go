package authentication

import (
	"errors"
	"net/http/httptest"
	"testing"
)

type ApiKeyAuthenticationValidatorMock struct {
	ApiKey string
	Calls  int
}

func (m *ApiKeyAuthenticationValidatorMock) GetApiKeyAuthenticationData(apiKey string) (ClaimMap, error) {
	m.ApiKey = apiKey
	m.Calls++

	if apiKey == "test" {
		claims := make(ClaimMap)
		claims[ClaimName] = &Claim{
			Name: ClaimName, Values: []string{"test_user"},
		}
		return claims, nil
	}
	if apiKey == "error" {
		return nil, errors.New("test error")
	}
	return nil, nil
}

func ApiKeyAuthenticationHandlerOptionMock(called *int) ApiKeyAuthenticationHandlerOption {
	return func(h *ApiKeyAuthenticationHandler) {
		*called++
	}
}

func Test_NewApiKeyAuthenticationHandler(t *testing.T) {
	optionCalled := 0
	validator := &ApiKeyAuthenticationValidatorMock{}
	handler := NewApiKeyAuthenticationHandler(validator, ApiKeyAuthenticationHandlerOptionMock(&optionCalled))
	internal := handler.(*ApiKeyAuthenticationHandler)
	if internal.validator != validator {
		t.Error("NewApiKeyAuthenticationHandler() failed: validator not set")
	}
	if internal.HeaderName != DefaultApiKeyHeaderName {
		t.Errorf("NewApiKeyAuthenticationHandler() failed: default header name not set: got %s, expected %s", internal.HeaderName, DefaultApiKeyHeaderName)
	}
	if internal.QueryParamName != DefaultApiKeyQueryParamName {
		t.Errorf("NewApiKeyAuthenticationHandler() failed: default query param name not set: got %s, expected %s", internal.QueryParamName, DefaultApiKeyQueryParamName)
	}
	if optionCalled != 1 {
		t.Error("NewApiKeyAuthenticationHandler() failed: options not set")
	}
}

func Test_ApiKeyAuthentication_HandleAuthentication_WithHeader(t *testing.T) {
	validator := &ApiKeyAuthenticationValidatorMock{}
	handler := &ApiKeyAuthenticationHandler{
		validator:      validator,
		HeaderName:     "api-key",
		QueryParamName: "api-key",
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("api-key", "test")
	context := handler.HandleAuthentication(request)

	if validator.ApiKey != "test" {
		t.Errorf("ApiKeyAuthenticationHandler.HandleAuthentication() failed: invalid api key: got %s, expected test", validator.ApiKey)
	}
	if validator.Calls != 1 {
		t.Errorf("ApiKeyAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 1", validator.Calls)
	}
	if context == nil {
		t.Error("ApiKeyAuthenticationHandler.HandleAuthentication() failed: nil context")
	} else {
		if !context.IsAuthenticated() {
			t.Error("ApiKeyAuthenticationHandler.HandleAuthentication() failed: not authenticated")
		}
		if context.GetName() != "test_user" {
			t.Errorf("ApiKeyAuthenticationHandler.HandleAuthentication() failed: invalid username: got %s, expected test_user", context.GetName())
		}
	}
}

func Test_ApiKeyAuthentication_HandleAuthentication_WithQueryParam(t *testing.T) {
	validator := &ApiKeyAuthenticationValidatorMock{}
	handler := &ApiKeyAuthenticationHandler{
		validator:      validator,
		HeaderName:     "api-key",
		QueryParamName: "api-key",
	}
	request := httptest.NewRequest("GET", "http://testing?api-key=test", nil)
	context := handler.HandleAuthentication(request)

	if validator.ApiKey != "test" {
		t.Errorf("ApiKeyAuthenticationHandler.HandleAuthentication() failed: invalid api key: got %s, expected test", validator.ApiKey)
	}
	if validator.Calls != 1 {
		t.Errorf("ApiKeyAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 1", validator.Calls)
	}
	if !context.IsAuthenticated() {
		t.Error("ApiKeyAuthenticationHandler.HandleAuthentication() failed: not authenticated")
	}
	if context.GetName() != "test_user" {
		t.Errorf("ApiKeyAuthenticationHandler.HandleAuthentication() failed: invalid username: got %s, expected test_user", context.GetName())
	}
}

func Test_ApiKeyAuthentication_HandleAuthentication_NilValidator(t *testing.T) {
	handler := &ApiKeyAuthenticationHandler{}
	request := httptest.NewRequest("GET", "http://testing", nil)
	context := handler.HandleAuthentication(request)
	if context != nil {
		t.Error("ApiKeyAuthenticationHandler.HandleAuthentication() failed: got valid context, expected nil")
	}
}

func Test_ApiKeyAuthentication_HandleAuthentication_ReturnsError(t *testing.T) {
	validator := &ApiKeyAuthenticationValidatorMock{}
	handler := &ApiKeyAuthenticationHandler{
		validator:      validator,
		HeaderName:     "api-key",
		QueryParamName: "api-key",
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("api-key", "error")
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("ApiKeyAuthenticationHandler.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 1 {
		t.Errorf("ApiKeyAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 1", validator.Calls)
	}
}

func Test_ApiKeyAuthentication_HandleAuthentication_ReturnsNil(t *testing.T) {
	validator := &ApiKeyAuthenticationValidatorMock{}
	handler := &ApiKeyAuthenticationHandler{
		validator:      validator,
		HeaderName:     "api-key",
		QueryParamName: "api-key",
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("api-key", "nil")
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("ApiKeyAuthenticationHandler.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 1 {
		t.Errorf("ApiKeyAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 1", validator.Calls)
	}
}
