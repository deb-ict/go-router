package authentication

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deb-ict/go-router"
)

type handlerMock struct {
	ReturnNilContext bool
}

func (h *handlerMock) HandleAuthentication(r *http.Request) Context {
	if h.ReturnNilContext {
		return nil
	}
	claims := make(ClaimMap)
	claims[ClaimName] = &Claim{
		Name: ClaimName, Values: []string{"test_user"},
	}
	claims["test"] = &Claim{
		Name: "test", Values: []string{"value"},
	}
	return &authenticationContext{
		authenticated: true,
		claims:        claims,
	}
}
func (h *handlerMock) EnsureDefaults() {
	h.ReturnNilContext = false
}

func MiddlewareOptionMock(called *int) MiddlewareOption {
	return func(m *Middleware) {
		*called++
	}
}

func Test_NewMiddleware(t *testing.T) {
	optionCalled := 0
	h := &handlerMock{}
	o := MiddlewareOptionMock(&optionCalled)
	m := NewMiddleware(h, o)

	if m == nil {
		t.Error("NewMiddleware failed: No instance")
	} else {
		if m.Handler == nil || m.Handler != h {
			t.Error("NewMiddleware failed: Handler not set")
		}
		if optionCalled != 1 {
			t.Error("NewMiddleware failed: Options not applied")
		}
	}
}

func Test_UseMiddleware(t *testing.T) {
	router := &router.Router{}
	handler := &handlerMock{}
	UseMiddleware(router, handler)

	if len(router.Middlewares()) != 1 {
		t.Error("UseMiddleware failed: Middleware not set on router")
	}
}

func Test_Middleware_WithNext(t *testing.T) {
	handler := &handlerMock{}
	middleware := &Middleware{
		Handler: handler,
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxValue := r.Context().Value(contextKey)
		if ctxValue == nil {
			t.Error("Middleware() failed: authentication context not set")
		} else {
			authContext := ctxValue.(Context)
			if authContext == nil {
				t.Error("Middleware() failed: authentication context not valid")
			} else {
				if !authContext.IsAuthenticated() {
					t.Error("Middleware() failed: authentication context not set as authenticated")
				}
				if authContext.GetName() != "test_user" {
					t.Errorf("Middleware() failed: authentication context name incorrect: got %s, expected test_user", authContext.GetName())
				}
			}
		}
	})
	test := middleware.Middleware(next)
	req := httptest.NewRequest("GET", "http://testing", nil)
	test.ServeHTTP(httptest.NewRecorder(), req)
}

func Test_Middleware_WithoutNext(t *testing.T) {
	handler := &handlerMock{}
	middleware := &Middleware{
		Handler: handler,
	}

	test := middleware.Middleware(nil)
	req := httptest.NewRequest("GET", "http://testing", nil)
	test.ServeHTTP(httptest.NewRecorder(), req)
}

func Test_Middleware_WithoutAuthContext(t *testing.T) {
	handler := &handlerMock{
		ReturnNilContext: true,
	}
	middleware := &Middleware{
		Handler: handler,
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxValue := r.Context().Value(contextKey)
		if ctxValue == nil {
			t.Error("Middleware() failed: authentication context not set")
		} else {
			authContext := ctxValue.(Context)
			if authContext == nil {
				t.Error("Middleware() failed: authentication context not valid")
			} else {
				if authContext.IsAuthenticated() {
					t.Error("Middleware() failed: authentication context set as authenticated")
				}
				if authContext.GetName() != "anonymous" {
					t.Errorf("Middleware() failed: authentication context name incorrect: got %s, expected anonymous", authContext.GetName())
				}
			}
		}
	})
	test := middleware.Middleware(next)
	req := httptest.NewRequest("GET", "http://testing", nil)
	test.ServeHTTP(httptest.NewRecorder(), req)
}
