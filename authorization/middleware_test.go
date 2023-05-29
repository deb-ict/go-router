package authorization

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deb-ict/go-router"
	"github.com/deb-ict/go-router/authentication"
)

func MiddlewareOptionMock(called *int) MiddlewareOption {
	return func(m *Middleware) {
		*called++
	}
}

func Test_NewMiddleware(t *testing.T) {
	optionCalled := 0
	o := MiddlewareOptionMock(&optionCalled)
	m := NewMiddleware(o)

	if m == nil {
		t.Error("NewMiddleware() failed: No instance")
	} else {
		if optionCalled != 1 {
			t.Error("NewMiddleware() failed: Options not applied")
		}
		if m.policies == nil {
			t.Error("NewMiddleware() failed: Policy collection not initialized")
		}
		if m.UnauthorizedHandler == nil {
			t.Error("NewMiddleware() failed: Default unauthorized handler not set")
		}
		if m.ForbiddenHandler == nil {
			t.Error("NewMiddleware() failed: Default forbidden handler not set")
		}
	}
}

func Test_UseMiddleware(t *testing.T) {
	router := &router.Router{}
	UseMiddleware(router)

	if len(router.Middlewares()) != 1 {
		t.Error("UseMiddleware() failed: Middleware not set on router")
	}
}

func Test_Middleware(t *testing.T) {
	var routeContextKey router.ContextKey = "router:route"

	unauthorizedCalled := 0
	unauthorized := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unauthorizedCalled++
	})

	forbiddenCalled := 0
	forbidden := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forbiddenCalled++
	})

	nextCalled := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled++
	})

	requirementCalled := 0
	requirement := &policyRequirementMock{
		result: true,
		calls:  &requirementCalled,
	}

	claims := make(authentication.ClaimMap, 0)
	claims.SetName("test")
	auth := authentication.NewContext(true, claims)

	route := &router.Route{}
	route.Authorize("test")
	middleware := &Middleware{
		policies:            make(map[string]Policy),
		UnauthorizedHandler: unauthorized,
		ForbiddenHandler:    forbidden,
	}
	middleware.policies["test"] = &policy{
		name:         "test",
		requirements: []Requirement{requirement},
	}
	test := middleware.Middleware(next)
	req := httptest.NewRequest("GET", "http://testing", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, routeContextKey, route)
	ctx = authentication.SetContext(ctx, auth)
	req = req.WithContext(ctx)
	test.ServeHTTP(httptest.NewRecorder(), req)

	if unauthorizedCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: Unauthorized called %x times, expected 0", unauthorizedCalled)
	}
	if forbiddenCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: Forbidden called %x times, expected 0", unauthorizedCalled)
	}
	if nextCalled != 1 {
		t.Errorf("Middleware(noRoute) failed: next called %x times, expected 1", unauthorizedCalled)
	}
	if requirementCalled != 1 {
		t.Errorf("Middleware(noRoute) failed: requirement called %x times, expected 1", requirementCalled)
	}
}

func Test_Middleware_RequirementFailed(t *testing.T) {
	var routeContextKey router.ContextKey = "router:route"

	unauthorizedCalled := 0
	unauthorized := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unauthorizedCalled++
	})

	forbiddenCalled := 0
	forbidden := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forbiddenCalled++
	})

	nextCalled := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled++
	})

	requirementCalled := 0
	requirement := &policyRequirementMock{
		result: false,
		calls:  &requirementCalled,
	}

	claims := make(authentication.ClaimMap, 0)
	claims.SetName("test")
	auth := authentication.NewContext(true, claims)

	route := &router.Route{}
	route.Authorize("test")
	middleware := &Middleware{
		policies:            make(map[string]Policy),
		UnauthorizedHandler: unauthorized,
		ForbiddenHandler:    forbidden,
	}
	middleware.policies["test"] = &policy{
		name:         "test",
		requirements: []Requirement{requirement},
	}
	test := middleware.Middleware(next)
	req := httptest.NewRequest("GET", "http://testing", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, routeContextKey, route)
	ctx = authentication.SetContext(ctx, auth)
	req = req.WithContext(ctx)
	test.ServeHTTP(httptest.NewRecorder(), req)

	if unauthorizedCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: Unauthorized called %x times, expected 0", unauthorizedCalled)
	}
	if forbiddenCalled != 1 {
		t.Errorf("Middleware(noRoute) failed: Forbidden called %x times, expected 1", unauthorizedCalled)
	}
	if nextCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: next called %x times, expected 0", unauthorizedCalled)
	}
	if requirementCalled != 1 {
		t.Errorf("Middleware(noRoute) failed: requirement called %x times, expected 1", requirementCalled)
	}
}

func Test_Middleware_PolicyNotFound(t *testing.T) {
	var routeContextKey router.ContextKey = "router:route"

	unauthorizedCalled := 0
	unauthorized := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unauthorizedCalled++
	})

	forbiddenCalled := 0
	forbidden := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forbiddenCalled++
	})

	nextCalled := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled++
	})

	claims := make(authentication.ClaimMap, 0)
	claims.SetName("test")
	auth := authentication.NewContext(true, claims)

	route := &router.Route{}
	route.Authorize("test")
	middleware := &Middleware{
		policies:            make(map[string]Policy),
		UnauthorizedHandler: unauthorized,
		ForbiddenHandler:    forbidden,
	}
	test := middleware.Middleware(next)
	req := httptest.NewRequest("GET", "http://testing", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, routeContextKey, route)
	ctx = authentication.SetContext(ctx, auth)
	req = req.WithContext(ctx)
	test.ServeHTTP(httptest.NewRecorder(), req)

	if unauthorizedCalled != 1 {
		t.Errorf("Middleware(noRoute) failed: Unauthorized called %x times, expected 1", unauthorizedCalled)
	}
	if forbiddenCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: Forbidden called %x times, expected 0", unauthorizedCalled)
	}
	if nextCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: next called %x times, expected 0", unauthorizedCalled)
	}
}

func Test_Middleware_Unauthorized(t *testing.T) {
	var routeContextKey router.ContextKey = "router:route"

	unauthorizedCalled := 0
	unauthorized := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unauthorizedCalled++
	})

	forbiddenCalled := 0
	forbidden := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forbiddenCalled++
	})

	nextCalled := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled++
	})

	route := &router.Route{}
	route.Authorize("test")
	middleware := &Middleware{
		policies:            make(map[string]Policy),
		UnauthorizedHandler: unauthorized,
		ForbiddenHandler:    forbidden,
	}
	test := middleware.Middleware(next)
	req := httptest.NewRequest("GET", "http://testing", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, routeContextKey, route)
	req = req.WithContext(ctx)
	test.ServeHTTP(httptest.NewRecorder(), req)

	if unauthorizedCalled != 1 {
		t.Errorf("Middleware(noRoute) failed: Unauthorized called %x times, expected 1", unauthorizedCalled)
	}
	if forbiddenCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: Forbidden called %x times, expected 0", unauthorizedCalled)
	}
	if nextCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: next called %x times, expected 0", unauthorizedCalled)
	}
}

func Test_Middleware_NoAuthorizationPolicy(t *testing.T) {
	var routeContextKey router.ContextKey = "router:route"

	unauthorizedCalled := 0
	unauthorized := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unauthorizedCalled++
	})

	forbiddenCalled := 0
	forbidden := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forbiddenCalled++
	})

	nextCalled := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled++
	})

	route := &router.Route{}
	middleware := &Middleware{
		policies:            make(map[string]Policy),
		UnauthorizedHandler: unauthorized,
		ForbiddenHandler:    forbidden,
	}
	test := middleware.Middleware(next)
	req := httptest.NewRequest("GET", "http://testing", nil)
	ctx := req.Context()
	ctx = context.WithValue(ctx, routeContextKey, route)
	req = req.WithContext(ctx)
	test.ServeHTTP(httptest.NewRecorder(), req)

	if unauthorizedCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: Unauthorized called %x times, expected 0", unauthorizedCalled)
	}
	if forbiddenCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: Forbidden called %x times, expected 0", unauthorizedCalled)
	}
	if nextCalled != 1 {
		t.Errorf("Middleware(noRoute) failed: next called %x times, expected 1", unauthorizedCalled)
	}
}

func Test_Middleware_NoRoute(t *testing.T) {
	unauthorizedCalled := 0
	unauthorized := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		unauthorizedCalled++
	})

	forbiddenCalled := 0
	forbidden := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forbiddenCalled++
	})

	nextCalled := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled++
	})

	middleware := &Middleware{
		policies:            make(map[string]Policy),
		UnauthorizedHandler: unauthorized,
		ForbiddenHandler:    forbidden,
	}
	test := middleware.Middleware(next)
	req := httptest.NewRequest("GET", "http://testing", nil)
	test.ServeHTTP(httptest.NewRecorder(), req)

	if unauthorizedCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: Unauthorized called %x times, expected 0", unauthorizedCalled)
	}
	if forbiddenCalled != 0 {
		t.Errorf("Middleware(noRoute) failed: Forbidden called %x times, expected 0", unauthorizedCalled)
	}
	if nextCalled != 1 {
		t.Errorf("Middleware(noRoute) failed: next called %x times, expected 1", unauthorizedCalled)
	}
}
