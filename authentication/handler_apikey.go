package authentication

import (
	"net/http"
)

type ApiKeyAuthenticationHandlerOption func(*ApiKeyAuthenticationHandler)

type ApiKeyAuthenticationValidator interface {
	GetApiKeyAuthenticationData(apiKey string) (ClaimMap, error)
}

type ApiKeyAuthenticationHandler struct {
	validator      ApiKeyAuthenticationValidator
	HeaderName     string
	QueryParamName string
}

func NewApiKeyAuthenticationHandler(validator ApiKeyAuthenticationValidator, opts ...ApiKeyAuthenticationHandlerOption) Handler {
	h := &ApiKeyAuthenticationHandler{
		validator: validator,
	}
	for _, opt := range opts {
		opt(h)
	}
	h.EnsureDefaults()

	return h
}

func (h *ApiKeyAuthenticationHandler) HandleAuthentication(r *http.Request) Context {
	if h.validator == nil {
		return nil
	}

	apiKey := r.Header.Get(h.HeaderName)
	if apiKey == "" {
		apiKey = r.URL.Query().Get(h.QueryParamName)
	}

	claims, err := h.validator.GetApiKeyAuthenticationData(apiKey)
	if err != nil {
		return nil
	}
	if claims == nil {
		return nil
	}

	return NewContext(true, claims)
}

func (h *ApiKeyAuthenticationHandler) EnsureDefaults() {
	if h.HeaderName == "" {
		h.HeaderName = DefaultApiKeyHeaderName
	}
	if h.QueryParamName == "" {
		h.QueryParamName = DefaultApiKeyQueryParamName
	}
}
