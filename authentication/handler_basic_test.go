package authentication

import (
	"encoding/base64"
	"errors"
	"net/http/httptest"
	"testing"
)

type BasicAuthenticationValidatorMock struct {
	Username string
	Password string
	Calls    int
}

func (m *BasicAuthenticationValidatorMock) GetUserAuthenticationData(username string, password string) (ClaimMap, error) {
	m.Username = username
	m.Password = password
	m.Calls++

	if username == "test_name" && password == "test_pass" {
		claims := make(ClaimMap)
		claims[ClaimName] = &Claim{
			Name: ClaimName, Values: []string{"test_user"},
		}
		return claims, nil
	}
	if username == "error" {
		return nil, errors.New("test error")
	}
	return nil, nil
}

func BasicAuthenticationHandlerOptionMock(called *int) BasicAuthenticationHandlerOption {
	return func(h *BasicAuthenticationHandler) {
		*called++
	}
}

func Test_NewBasicAuthenticationHandler(t *testing.T) {
	optionCalled := 0
	validator := &BasicAuthenticationValidatorMock{}
	handler := NewBasicAuthenticationHandler(validator, BasicAuthenticationHandlerOptionMock(&optionCalled))
	internal := handler.(*BasicAuthenticationHandler)
	if internal.validator != validator {
		t.Error("NewBasicAuthenticationHandler() failed: validator not set")
	}
	if optionCalled != 1 {
		t.Error("NewBasicAuthenticationHandler() failed: options not set")
	}
}

func Test_BasicAuthentication_HandleAuthentication(t *testing.T) {
	credentials := "test_name:test_pass"
	token := base64.StdEncoding.EncodeToString([]byte(credentials))

	validator := BasicAuthenticationValidatorMock{}
	handler := &BasicAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Basic "+token)
	context := handler.HandleAuthentication(request)

	if validator.Username != "test_name" {
		t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid username: got %s, expected test_name", validator.Username)
	}
	if validator.Password != "test_pass" {
		t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid password: got %s, expected test_pass", validator.Password)
	}
	if validator.Calls != 1 {
		t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 1", validator.Calls)
	}
	if context == nil {
		t.Error("BasicAuthenticationHandler.HandleAuthentication() failed: nil context")
	} else {
		if !context.IsAuthenticated() {
			t.Error("BasicAuthenticationHandler.HandleAuthentication() failed: not authenticated")
		}
		if context.GetName() != "test_user" {
			t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid username: got %s, expected test_user", context.GetName())
		}
	}
}

func Test_BasicAuthentication_HandleAuthentication_NilValidator(t *testing.T) {
	handler := &BasicAuthenticationHandler{}
	request := httptest.NewRequest("GET", "http://testing", nil)
	context := handler.HandleAuthentication(request)
	if context != nil {
		t.Error("BasicAuthenticationHandler.HandleAuthentication() failed: got valid context, expected nil")
	}
}

func Test_BasicAuthentication_HandleAuthentication_HeaderNotSet(t *testing.T) {
	validator := BasicAuthenticationValidatorMock{}
	handler := &BasicAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BasicAuthenticationHandler.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 0 {
		t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 0", validator.Calls)
	}
}

func Test_BasicAuthentication_HandleAuthentication_InvalidHeaderValuePrefix(t *testing.T) {
	validator := BasicAuthenticationValidatorMock{}
	handler := &BasicAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Invalid test")
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BasicAuthenticationHandler.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 0 {
		t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 0", validator.Calls)
	}
}

func Test_BasicAuthentication_HandleAuthentication_ZeroLengthToken(t *testing.T) {
	validator := BasicAuthenticationValidatorMock{}
	handler := &BasicAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Basic ")
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BasicAuthenticationHandler.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 0 {
		t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 0", validator.Calls)
	}
}

func Test_BasicAuthentication_HandleAuthentication_NoSeperator(t *testing.T) {
	credentials := "test_name_test_pass"
	token := base64.StdEncoding.EncodeToString([]byte(credentials))

	validator := BasicAuthenticationValidatorMock{}
	handler := &BasicAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Basic "+token)
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BasicAuthenticationHandler.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 0 {
		t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 0", validator.Calls)
	}
}

func Test_BasicAuthentication_HandleAuthentication_NoBase64(t *testing.T) {
	validator := BasicAuthenticationValidatorMock{}
	handler := &BasicAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Basic text_token")
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BasicAuthenticationHandler.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 0 {
		t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 0", validator.Calls)
	}
}

//TODO: Test header not set
//TODO: Test incorrect header prefix
//TODO: Test zero lenght token
//TODO: Test invalid token (no :)

//TODO: Test validator returns error
//TODO: Test validator returns nil
//TODO: Test validator returns context

func Test_BasicAuthentication_HandleAuthentication_ReturnsError(t *testing.T) {
	credentials := "error:error"
	token := base64.StdEncoding.EncodeToString([]byte(credentials))

	validator := BasicAuthenticationValidatorMock{}
	handler := &BasicAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Basic "+token)
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BasicAuthenticationHandler.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 1 {
		t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 1", validator.Calls)
	}
}

func Test_BasicAuthentication_HandleAuthentication_ReturnsNil(t *testing.T) {
	credentials := "unknown:unknown"
	token := base64.StdEncoding.EncodeToString([]byte(credentials))

	validator := BasicAuthenticationValidatorMock{}
	handler := &BasicAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Basic "+token)
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BasicAuthenticationHandler.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 1 {
		t.Errorf("BasicAuthenticationHandler.HandleAuthentication() failed: invalid number of calls: got %d, expected 1", validator.Calls)
	}
}
