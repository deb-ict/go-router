package authentication

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type BasicAuthenticationHandlerOption func(*BasicAuthenticationHandler)

type BasicAuthenticationValidator interface {
	GetUserAuthenticationData(username string, password string) (ClaimMap, error)
}

type BasicAuthenticationHandler struct {
	validator BasicAuthenticationValidator
}

func NewBasicAuthenticationHandler(validator BasicAuthenticationValidator, opts ...BasicAuthenticationHandlerOption) Handler {
	h := &BasicAuthenticationHandler{
		validator: validator,
	}
	for _, opt := range opts {
		opt(h)
	}
	h.EnsureDefaults()

	return h
}

func (h *BasicAuthenticationHandler) HandleAuthentication(r *http.Request) Context {
	if h.validator == nil {
		return nil
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil
	}
	if len(auth) < len(BasicTokenPrex) || !equalFold(auth[:len(BasicTokenPrex)], BasicTokenPrex) {
		return nil
	}
	token := auth[len(BasicTokenPrex):]

	tokenData, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil
	}
	tokenString := string(tokenData)
	username, password, ok := strings.Cut(tokenString, ":")
	if !ok {
		return nil
	}

	claims, err := h.validator.GetUserAuthenticationData(username, password)
	if err != nil {
		return nil
	}
	if claims == nil {
		return nil
	}

	return NewContext(true, claims)
}

func (h *BasicAuthenticationHandler) EnsureDefaults() {

}
