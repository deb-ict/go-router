package authentication

import (
	"net/http"
)

type ApiKeyHandlerOption func(*ApiKeyHandler)

type ApiKeyValidator interface {
	GetApiKeyAuthenticationData(apiKey string) (ClaimMap, error)
}

type ApiKeyHandler struct {
	validator      ApiKeyValidator
	HeaderName     string
	QueryParamName string
}

func NewApiKeyHandler(validator ApiKeyValidator, opts ...ApiKeyHandlerOption) Handler {
	h := &ApiKeyHandler{
		validator: validator,
	}
	for _, opt := range opts {
		opt(h)
	}
	h.EnsureDefaults()

	return h
}

func (h *ApiKeyHandler) HandleAuthentication(r *http.Request) Context {
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

func (h *ApiKeyHandler) EnsureDefaults() {
	if h.HeaderName == "" {
		h.HeaderName = DefaultApiKeyHeaderName
	}
	if h.QueryParamName == "" {
		h.QueryParamName = DefaultApiKeyQueryParamName
	}
}
