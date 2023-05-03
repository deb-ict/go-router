package authentication

import (
	"net/http"
)

type BearerTokenHandlerOption func(*BearerTokenHandler)

type BearerTokenValidator interface {
	GetBearerAuthenticationData(token string) (ClaimMap, error)
}

type BearerTokenHandler struct {
	validator BearerTokenValidator
}

func NewBearerHandler(validator BearerTokenValidator, opts ...BearerTokenHandlerOption) Handler {
	h := &BearerTokenHandler{
		validator: validator,
	}
	for _, opt := range opts {
		opt(h)
	}
	h.EnsureDefaults()

	return h
}

func (h *BearerTokenHandler) HandleAuthentication(r *http.Request) Context {
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

	claims, err := h.validator.GetBearerAuthenticationData(token)
	if err != nil {
		return nil
	}
	if claims == nil {
		return nil
	}

	return NewContext(true, claims)
}

func (h *BearerTokenHandler) EnsureDefaults() {

}
