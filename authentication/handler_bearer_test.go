package authentication

import (
	"errors"
	"net/http/httptest"
	"testing"
)

type BearerAuthenticationValidatorMock struct {
	Token string
	Calls int
}

func (m *BearerAuthenticationValidatorMock) GetBearerAuthenticationData(token string) (ClaimMap, error) {
	m.Token = token
	m.Calls++

	if token == "test" {
		claims := make(ClaimMap)
		claims[ClaimName] = &Claim{
			Name: ClaimName, Values: []string{"test_user"},
		}
		return claims, nil
	}
	if token == "error" {
		return nil, errors.New("test error")
	}

	return nil, nil
}

func BearerAuthenticationHandlerOptionMock(called *int) BearerAuthenticationHandlerOption {
	return func(h *BearerAuthenticationHandler) {
		*called++
	}
}

func Test_NewBearerAuthenticationHandler(t *testing.T) {
	optionCalled := 0
	validator := &BearerAuthenticationValidatorMock{}
	handler := NewBearerAuthenticationHandler(validator, BearerAuthenticationHandlerOptionMock(&optionCalled))
	internal := handler.(*BearerAuthenticationHandler)
	if internal.validator != validator {
		t.Error("NewBearerAuthenticationHandler() failed: validator not set")
	}
	if optionCalled != 1 {
		t.Error("NewBearerAuthenticationHandler() failed: options not set")
	}
}

func Test_BearerAuthentication_HandleAuthentication(t *testing.T) {
	validator := BearerAuthenticationValidatorMock{}
	handler := &BearerAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Bearer test")
	context := handler.HandleAuthentication(request)

	if validator.Token != "test" {
		t.Errorf("BearerAuthentication.HandleAuthentication() failed: invalid token: got %s, expected test", validator.Token)
	}
	if validator.Calls != 1 {
		t.Errorf("BearerAuthentication.HandleAuthentication() failed: invalid number of calls: got %d, expected 1", validator.Calls)
	}
	if context == nil {
		t.Error("BearerAuthentication.HandleAuthentication() failed: nil context")
	} else {
		if !context.IsAuthenticated() {
			t.Error("BearerAuthentication.HandleAuthentication() failed: not authenticated")
		}
		if context.GetName() != "test_user" {
			t.Errorf("BearerAuthentication.HandleAuthentication() failed: invalid username: got %s, expected test_user", context.GetName())
		}
	}
}

func Test_BearerAuthentication_HandleAuthentication_NilValidator(t *testing.T) {
	handler := &BearerAuthenticationHandler{}
	request := httptest.NewRequest("GET", "http://testing", nil)
	context := handler.HandleAuthentication(request)
	if context != nil {
		t.Error("BearerAuthenticationHandler.HandleAuthentication() failed: got valid context, expected nil")
	}
}

func Test_BearerAuthentication_HandleAuthentication_HeaderNotSet(t *testing.T) {
	validator := BearerAuthenticationValidatorMock{}
	handler := &BearerAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BearerAuthentication.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 0 {
		t.Errorf("BearerAuthentication.HandleAuthentication() failed: invalid number of calls: got %d, expected 0", validator.Calls)
	}
}

func Test_BearerAuthentication_HandleAuthentication_InvalidHeaderValuePrefix(t *testing.T) {
	validator := BearerAuthenticationValidatorMock{}
	handler := &BearerAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Invalid test")
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BearerAuthentication.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 0 {
		t.Errorf("BearerAuthentication.HandleAuthentication() failed: invalid number of calls: got %d, expected 0", validator.Calls)
	}
}

func Test_BearerAuthentication_HandleAuthentication_ZeroLengthToken(t *testing.T) {
	validator := BearerAuthenticationValidatorMock{}
	handler := &BearerAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Bearer ")
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BearerAuthentication.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 0 {
		t.Errorf("BearerAuthentication.HandleAuthentication() failed: invalid number of calls: got %d, expected 0", validator.Calls)
	}
}

func Test_BearerAuthentication_HandleAuthentication_ReturnsError(t *testing.T) {
	validator := BearerAuthenticationValidatorMock{}
	handler := &BearerAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Bearer error")
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BearerAuthentication.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 1 {
		t.Errorf("BearerAuthentication.HandleAuthentication() failed: invalid number of calls: got %d, expected 1", validator.Calls)
	}
}

func Test_BearerAuthentication_HandleAuthentication_ReturnsNil(t *testing.T) {
	validator := BearerAuthenticationValidatorMock{}
	handler := &BearerAuthenticationHandler{
		validator: &validator,
	}
	request := httptest.NewRequest("GET", "http://testing", nil)
	request.Header.Add("Authorization", "Bearer nil")
	context := handler.HandleAuthentication(request)

	if context != nil {
		t.Error("BearerAuthentication.HandleAuthentication() failed: get valid context, expected nil")
	}
	if validator.Calls != 1 {
		t.Errorf("BearerAuthentication.HandleAuthentication() failed: invalid number of calls: got %d, expected 1", validator.Calls)
	}
}
