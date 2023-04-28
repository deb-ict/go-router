package authentication

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type BasicAuthenticationHandlerOption func(*BasicAuthenticationHandler)

type BasicAuthenticationHandler struct {
}

func NewBasicAuthenticationHandler(opts ...BasicAuthenticationHandlerOption) AuthenticationHandler {
	h := &BasicAuthenticationHandler{}
	for _, opt := range opts {
		opt(h)
	}
	h.EnsureDefaults()

	return h
}

func (h *BasicAuthenticationHandler) HandleAuthentication(r *http.Request) *AuthenticationContext {
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

	if username != "my-user" || password != "my-pass" {
		return &AuthenticationContext{}
	}

	return nil
}

func (h *BasicAuthenticationHandler) EnsureDefaults() {

}
