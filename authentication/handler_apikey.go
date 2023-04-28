package authentication

import (
	"net/http"
)

type ApiKeyAuthenticationHandlerOption func(*ApiKeyAuthenticationHandler)

type ApiKeyAuthenticationHandler struct {
	HeaderName     string
	QueryParamName string
}

func NewApiKeyAuthenticationHandler(opts ...ApiKeyAuthenticationHandlerOption) AuthenticationHandler {
	h := &ApiKeyAuthenticationHandler{}
	for _, opt := range opts {
		opt(h)
	}
	h.EnsureDefaults()

	return h
}

func (h *ApiKeyAuthenticationHandler) HandleAuthentication(r *http.Request) *AuthenticationContext {
	apiKey := r.Header.Get(h.HeaderName)
	if apiKey == "" {
		apiKey = r.URL.Query().Get(h.QueryParamName)
	}

	if apiKey != "my-api-key" {
		return &AuthenticationContext{}
	}

	return nil
}

func (h *ApiKeyAuthenticationHandler) EnsureDefaults() {
	if h.HeaderName == "" {
		h.HeaderName = DefaultApiKeyHeaderName
	}
	if h.QueryParamName == "" {
		h.QueryParamName = DefaultApiKeyQueryParamName
	}
}
