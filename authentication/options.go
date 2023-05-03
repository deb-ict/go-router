package authentication

func WithApiKeyAuthenticationHeaderName(name string) ApiKeyHandlerOption {
	return func(h *ApiKeyHandler) {
		h.HeaderName = name
	}
}

func WithApiKeyAuthenticationQueryParamName(name string) ApiKeyHandlerOption {
	return func(h *ApiKeyHandler) {
		h.QueryParamName = name
	}
}
