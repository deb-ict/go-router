package authentication

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
