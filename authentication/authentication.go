package authentication

import (
	"context"
	"strings"

	"github.com/deb-ict/go-router"
)

const (
	DefaultApiKeyHeaderName     string = "X-API-KEY"
	DefaultApiKeyQueryParamName string = "api_key"
	AuthorizationHeaderName     string = "Authorization"
	BearerTokenPrefix           string = "Bearer "
	BasicTokenPrex              string = "Basic "
)

const (
	authenticationContextKey router.ContextKey = "router::authentication"
)

type Claim struct {
	Name   string
	Values []string
}

type ClaimMap map[string]Claim

type AuthenticationContext interface {
	IsAuthenticated() bool
	GetClaimMap() ClaimMap
	GetClaim(name string) Claim
	GetClaimValue(name string, index int) string
	HasClaim(name string) bool
	HasClaimValue(name string, value string) bool
}

type authenticationContext struct {
	authenticated bool
	claims        map[string]Claim
}

func GetAuthenticationContext(ctx context.Context) AuthenticationContext {
	value := ctx.Value(authenticationContextKey)
	if value == nil {
		return nil
	}
	return value.(AuthenticationContext)
}

func SetAuthenticationContext(ctx context.Context, auth AuthenticationContext) context.Context {
	return context.WithValue(ctx, authenticationContextKey, auth)
}

func (ctx *authenticationContext) IsAuthenticated() bool {
	return ctx.authenticated
}

func (ctx *authenticationContext) GetClaimMap() ClaimMap {
	return ctx.claims
}

func (ctx *authenticationContext) GetClaim(name string) Claim {
	return ctx.claims[name]
}

func (ctx *authenticationContext) GetClaimValue(name string, index int) string {
	c := ctx.GetClaim(name)
	if len(c.Values) <= index {
		return ""
	}
	return c.Values[index]
}

func (ctx *authenticationContext) HasClaim(name string) bool {
	_, ok := ctx.claims[name]
	return ok
}

func (ctx *authenticationContext) HasClaimValue(name string, value string) bool {
	c := ctx.GetClaim(name)
	return c.HasValue(value)
}

func (c *Claim) First() string {
	if len(c.Values) == 0 {
		return ""
	}
	return c.Values[0]
}

func (c *Claim) HasValue(value string) bool {
	for _, v := range c.Values {
		if strings.EqualFold(v, value) {
			return true
		}
	}
	return false
}
