package authentication

import (
	"context"

	"github.com/deb-ict/go-router"
)

const (
	DefaultApiKeyHeaderName     string = "X-Api-Key"
	DefaultApiKeyQueryParamName string = "api-key"
	AuthorizationHeaderName     string = "Authorization"
	BearerTokenPrefix           string = "Bearer "
	BasicTokenPrex              string = "Basic "
)

const (
	authenticationContextKey router.ContextKey = "router::authentication"
)

type AuthenticationContext struct {
}

func GetAuthenticationContext(ctx context.Context) *AuthenticationContext {
	value := ctx.Value(authenticationContextKey)
	if value == nil {
		return nil
	}
	return value.(*AuthenticationContext)
}

func SetAuthenticationContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, authenticationContextKey, nil)
}
