package authentication

import (
	"context"

	"github.com/deb-ict/go-router"
)

const (
	contextKey router.ContextKey = "router::authentication"
)

type Context interface {
	IsAuthenticated() bool
	GetSubjectId() string
	GetName() string
	GetRoles() []string
	GetScopes() []string
	GetClaim(name string) *Claim
	GetClaimValue(name string, index int) string
	HasClaim(name string) bool
	HasClaimValue(name string, value string) bool
	HasRole(name string) bool
	HasScope(name string) bool
}

type authenticationContext struct {
	authenticated bool
	claims        ClaimMap
}

func NewContext(authenticated bool, claims ClaimMap) Context {
	if claims == nil {
		claims = make(ClaimMap)
	}
	return &authenticationContext{
		authenticated: authenticated,
		claims:        claims,
	}
}

func AnonymouseContext() Context {
	claims := make(ClaimMap)
	claims.SetName("anonymous")
	return NewContext(false, claims)
}

func GetContext(ctx context.Context) Context {
	value := ctx.Value(contextKey)
	if value == nil {
		return AnonymouseContext()
	}
	return value.(Context)
}

func SetContext(ctx context.Context, value Context) context.Context {
	return context.WithValue(ctx, contextKey, value)
}

func (ctx *authenticationContext) IsAuthenticated() bool {
	return ctx.authenticated
}

func (ctx *authenticationContext) GetSubjectId() string {
	claim := ctx.GetClaim(ClaimSubjectId)
	return claim.First()
}

func (ctx *authenticationContext) GetName() string {
	claim := ctx.GetClaim(ClaimName)
	return claim.First()
}

func (ctx *authenticationContext) GetRoles() []string {
	claim := ctx.GetClaim(ClaimRole)
	return claim.Values
}

func (ctx *authenticationContext) GetScopes() []string {
	claim := ctx.GetClaim(ClaimScope)
	return claim.Values
}

func (ctx *authenticationContext) GetClaim(name string) *Claim {
	return ctx.claims.GetClaim(name)
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

func (ctx *authenticationContext) HasRole(name string) bool {
	return ctx.HasClaimValue(ClaimRole, name)
}

func (ctx *authenticationContext) HasScope(name string) bool {
	return ctx.HasClaimValue(ClaimScope, name)
}
