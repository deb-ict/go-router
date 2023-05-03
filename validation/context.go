package validation

import (
	"context"

	"github.com/deb-ict/go-router"
)

const (
	contextKey router.ContextKey = "router::validation"
)

type ErrorValues []string
type ErrorMap map[string]ErrorValues

type Context interface {
	Reset(name string)
	SetMessage(message string)
	GetMessage() string
	AddError(name string, message string)
	GetErrors() ErrorMap
	HasErrors() bool
}

type validationContext struct {
	message string
	errors  ErrorMap
}

func NewContext() Context {
	return &validationContext{
		errors: make(ErrorMap),
	}
}

func GetContext(ctx context.Context) Context {
	value := ctx.Value(contextKey)
	if value == nil {
		return NewContext()
	}
	return value.(Context)
}

func SetContext(ctx context.Context, value Context) context.Context {
	return context.WithValue(ctx, contextKey, value)
}

func (ctx *validationContext) Reset(name string) {
	ctx.errors[name] = make([]string, 0)
}

func (ctx *validationContext) SetMessage(message string) {
	ctx.message = message
}

func (ctx *validationContext) GetMessage() string {
	return ctx.message
}

func (ctx *validationContext) AddError(name string, message string) {
	if ctx.errors == nil {
		ctx.errors = make(ErrorMap)
	}
	if _, ok := ctx.errors[name]; !ok {
		ctx.errors[name] = make([]string, 0)
	}
	ctx.errors[name] = append(ctx.errors[name], message)
}

func (ctx *validationContext) GetErrors() ErrorMap {
	if ctx.errors == nil {
		return make(ErrorMap)
	}
	return ctx.errors
}

func (ctx *validationContext) HasErrors() bool {
	if ctx.errors == nil {
		return false
	}
	return len(ctx.errors) > 0
}
