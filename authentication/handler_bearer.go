package authentication

import (
	"net/http"
)

type BearerAuthenticationHandlerOption func(*BearerAuthenticationHandler)

type BearerAuthenticationHandler struct {
}

func NewBearerAuthenticationHandler(opts ...BearerAuthenticationHandlerOption) AuthenticationHandler {
	h := &BearerAuthenticationHandler{}
	for _, opt := range opts {
		opt(h)
	}
	h.EnsureDefaults()

	return h
}

func (h *BearerAuthenticationHandler) HandleAuthentication(r *http.Request) *AuthenticationContext {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil
	}
	if len(auth) < len(BearerTokenPrefix) || !equalFold(auth[:len(BearerTokenPrefix)], BearerTokenPrefix) {
		return nil
	}
	token := auth[len(BearerTokenPrefix):]

	if token != "my-token" {
		return &AuthenticationContext{}
	}

	return nil
}

func (h *BearerAuthenticationHandler) EnsureDefaults() {

}
