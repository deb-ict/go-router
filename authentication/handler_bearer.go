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

func (h *BearerAuthenticationHandler) HandleAuthentication(r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return
	}
	if len(auth) < len(BearerTokenPrefix) || !equalFold(auth[:len(BearerTokenPrefix)], BearerTokenPrefix) {
		return
	}
	token := auth[len(BearerTokenPrefix):]

	if token != "my-token" {
		return
	}
}

func (h *BearerAuthenticationHandler) EnsureDefaults() {

}
