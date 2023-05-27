package authentication

import (
	"net/http"
)

type BearerAuthenticationHandlerOption func(*BearerAuthenticationHandler)

type BearerAuthenticationValidator interface {
	GetBearerAuthenticationData(token string) (ClaimMap, error)
}

type BearerAuthenticationHandler struct {
	validator BearerAuthenticationValidator
}

func NewBearerAuthenticationHandler(validator BearerAuthenticationValidator, opts ...BearerAuthenticationHandlerOption) Handler {
	h := &BearerAuthenticationHandler{
		validator: validator,
	}
	for _, opt := range opts {
		opt(h)
	}
	h.EnsureDefaults()

	return h
}

func (h *BearerAuthenticationHandler) HandleAuthentication(r *http.Request) Context {
	if h.validator == nil {
		return nil
	}

	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil
	}
	if len(auth) < len(BearerTokenPrefix) || !equalFold(auth[:len(BearerTokenPrefix)], BearerTokenPrefix) {
		return nil
	}
	token := auth[len(BearerTokenPrefix):]
	if token == "" {
		return nil
	}

	claims, err := h.validator.GetBearerAuthenticationData(token)
	if err != nil {
		return nil
	}
	if claims == nil {
		return nil
	}

	return NewContext(true, claims)
}

func (h *BearerAuthenticationHandler) EnsureDefaults() {

}
