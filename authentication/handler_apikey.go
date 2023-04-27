package authentication

import (
	"net/http"
)

type ApiKeyAuthenticationHandlerOption func(*ApiKeyAuthenticationHandler)

type ApiKeyAuthenticationHandler struct {
	HeaderName     string
	QueryParamName string
}

func WithApiKeyAuthenticationHeaderName(name string) ApiKeyAuthenticationHandlerOption {
	return func(h *ApiKeyAuthenticationHandler) {
		h.HeaderName = name
	}
}

func WithApiKeyAuthenticationQueryParamName(name string) ApiKeyAuthenticationHandlerOption {
	return func(h *ApiKeyAuthenticationHandler) {
		h.QueryParamName = name
	}
}

func NewApiKeyAuthenticationHandler(opts ...ApiKeyAuthenticationHandlerOption) AuthenticationHandler {
	h := &ApiKeyAuthenticationHandler{}
	for _, opt := range opts {
		opt(h)
	}
	h.EnsureDefaults()

	return h
}

func (h *ApiKeyAuthenticationHandler) HandleAuthentication(r *http.Request) {
	apiKey := r.Header.Get(h.HeaderName)
	if apiKey == "" {
		apiKey = r.URL.Query().Get(h.QueryParamName)
	}

	if apiKey != "my-api-key" {
		return
	}
}

func (h *ApiKeyAuthenticationHandler) EnsureDefaults() {
	if h.HeaderName == "" {
		h.HeaderName = DefaultApiKeyHeaderName
	}
	if h.QueryParamName == "" {
		h.QueryParamName = DefaultApiKeyQueryParamName
	}
}
