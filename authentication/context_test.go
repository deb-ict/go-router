package authentication

import (
	"context"
	"testing"
)

func Test_NewContext_SetsClaims(t *testing.T) {
	claims := make(ClaimMap)
	claims["test"] = &Claim{
		Name:   "test",
		Values: []string{"value"},
	}
	context := NewContext(true, claims)
	internal := context.(*authenticationContext)

	if !internal.authenticated {
		t.Error("NewContext() failed: IsAuthentication not set")
	}
	if len(internal.claims) != 1 {
		t.Error("NewContext() failed: Initial claim value not found")
	}
}

func Test_NewContext_CreatesEmptyClaims(t *testing.T) {
	context := NewContext(true, nil)
	internal := context.(*authenticationContext)

	if !internal.authenticated {
		t.Error("NewContext() failed: IsAuthentication not set")
	}
	if internal.claims == nil {
		t.Error("NewContext() failed: Claims not initialized")
	}
}

func Test_AnonymousContext(t *testing.T) {
	context := AnonymouseContext()
	internal := context.(*authenticationContext)

	if internal.authenticated {
		t.Error("AnonymousContext() failed: IsAuthentication is set true")
	}
	c, ok := internal.claims[ClaimName]
	if !ok {
		t.Error("AnonymousContext() failed: Name claim not set")
	} else if c.First() != "anonymous" {
		t.Error("AnonymousContext() failed: Name claim value must be anonymous")
	}
}

func Test_GetContext_ReturnsAnonymousContextWhenNotSet(t *testing.T) {
	ctx := context.TODO()
	context := GetContext(ctx)
	if context == nil {
		t.Error("GetContext() failed: No default context created")
	}
	if context.IsAuthenticated() {
		t.Error("GetContext() failed: Default context set as authenticated")
	}
	if context.GetName() != "anonymous" {
		t.Error("GetContext() failed: Default context name not set to anonymous")
	}
}

func Test_GetContext_ReturnsContext(t *testing.T) {
	claims := make(ClaimMap)
	claims["test"] = &Claim{
		Name:   "test",
		Values: []string{"value"},
	}
	expected := NewContext(true, claims)
	ctx := context.TODO()
	ctx = context.WithValue(ctx, contextKey, expected)

	result := GetContext(ctx)
	if result != expected {
		t.Error("SetContext() failed: authentication context not set")
	}
}

func Test_SetContext(t *testing.T) {
	claims := make(ClaimMap)
	claims["test"] = &Claim{
		Name:   "test",
		Values: []string{"value"},
	}
	expected := NewContext(true, claims)
	ctx := context.TODO()
	ctx = SetContext(ctx, expected)

	result := ctx.Value(contextKey)
	if result != expected {
		t.Error("SetContext() failed: authentication context not set")
	}
}

func Test_AuthenticationContext_IsAuthenticated(t *testing.T) {
	type testCase struct {
		ctx      *authenticationContext
		expected bool
	}
	tests := []testCase{
		{&authenticationContext{authenticated: true}, true},
		{&authenticationContext{authenticated: false}, false},
	}

	for _, tc := range tests {
		result := tc.ctx.IsAuthenticated()
		if result != tc.expected {
			t.Errorf("AuthenticationContext.IsAuthenticated() failed: got %v, expected %v", result, tc.expected)
		}
	}
}

func Test_AuthenticationContext_GetSubjectId(t *testing.T) {
	expected := "123"
	ctx := &authenticationContext{
		authenticated: true,
		claims: ClaimMap{ClaimSubjectId: &Claim{
			ClaimSubjectId, []string{expected},
		}},
	}

	result := ctx.GetSubjectId()
	if result != expected {
		t.Errorf("AuthenticationContext.GetSubjectId() failed: got %v, expected %v", result, expected)
	}
}

func Test_AuthenticationContext_GetName(t *testing.T) {
	expected := "tester"
	ctx := &authenticationContext{
		authenticated: true,
		claims: ClaimMap{ClaimName: &Claim{
			ClaimName, []string{expected},
		}},
	}

	result := ctx.GetName()
	if result != expected {
		t.Errorf("AuthenticationContext.GetName() failed: got %v, expected %v", result, expected)
	}
}

func Test_AuthenticationContext_GetRoles(t *testing.T) {
	expected := []string{"role1", "role2"}
	ctx := &authenticationContext{
		authenticated: true,
		claims: ClaimMap{ClaimRole: &Claim{
			ClaimRole, expected,
		}},
	}

	result := ctx.GetRoles()
	if len(result) != len(expected) {
		t.Fatalf("AuthenticationContext.GetRoles() failed: incorrect number of values: got %v, expected %v", len(result), len(expected))
	}
	for i := 0; i < len(result); i++ {
		if result[i] != expected[i] {
			t.Errorf("AuthenticationContext.Roles() failed: incorrect values at index %v: got %v, expected %v", i, result[i], expected[i])
		}
	}
}

func Test_AuthenticationContext_GetScopes(t *testing.T) {
	expected := []string{"scope1", "scope2"}
	ctx := &authenticationContext{
		authenticated: true,
		claims: ClaimMap{ClaimScope: &Claim{
			ClaimScope, expected,
		}},
	}

	result := ctx.GetScopes()
	if len(result) != len(expected) {
		t.Fatalf("AuthenticationContext.GetScopes() failed: incorrect number of values: got %v, expected %v", len(result), len(expected))
	}
	for i := 0; i < len(result); i++ {
		if result[i] != expected[i] {
			t.Errorf("AuthenticationContext.GetScopes() failed: incorrect values at index %v: got %v, expected %v", i, result[i], expected[i])
		}
	}
}

func Test_AuthenticationContext_GetClaim(t *testing.T) {
	ctx := &authenticationContext{
		authenticated: true,
		claims: ClaimMap{
			"test1": &Claim{"test1", []string{"value1", "value2"}},
			"test2": &Claim{"test2", []string{}},
		},
	}
	result := ctx.GetClaim("test1")
	compareClaims(t, "AuthenticationContext.GetClaim()", result, ctx.claims["test1"])
}

func Test_AuthenticationContext_GetClaimValue(t *testing.T) {
	type testCase struct {
		name     string
		index    int
		expected string
	}
	tests := []testCase{
		{"test1", 0, "value1"},
		{"test1", 1, "value2"},
		{"test1", 2, ""},
		{"test2", 0, ""},
		{"test3", 0, ""},
	}
	ctx := &authenticationContext{
		authenticated: true,
		claims: ClaimMap{
			"test1": &Claim{"test1", []string{"value1", "value2"}},
			"test2": &Claim{"test2", []string{}},
		},
	}

	for _, tc := range tests {
		result := ctx.GetClaimValue(tc.name, tc.index)
		if result != tc.expected {
			t.Errorf("AuthenticationContext.GetClaimValue(%s, %d) failed: got %v, expected %v", tc.name, tc.index, result, tc.expected)
		}
	}
}

func Test_AuthenticationContext_HasClaim(t *testing.T) {
	type testCase struct {
		name     string
		expected bool
	}
	tests := []testCase{
		{"test1", true},
		{"test2", true},
		{"test3", false},
	}
	ctx := &authenticationContext{
		authenticated: true,
		claims: ClaimMap{
			"test1": &Claim{"test1", []string{"value1", "value2"}},
			"test2": &Claim{"test2", []string{}},
		},
	}

	for _, tc := range tests {
		result := ctx.HasClaim(tc.name)
		if result != tc.expected {
			t.Errorf("AuthenticationContext.HasClaim(%s) failed: got %v, expected %v", tc.name, result, tc.expected)
		}
	}
}

func Test_AuthenticationContext_HasClaimValue(t *testing.T) {
	type testCase struct {
		name     string
		value    string
		expected bool
	}
	tests := []testCase{
		{"test1", "value1", true},
		{"test1", "value3", false},
		{"test2", "value1", false},
		{"test3", "value1", false},
	}
	ctx := &authenticationContext{
		authenticated: true,
		claims: ClaimMap{
			"test1": &Claim{"test1", []string{"value1", "value2"}},
			"test2": &Claim{"test2", []string{}},
		},
	}

	for _, tc := range tests {
		result := ctx.HasClaimValue(tc.name, tc.value)
		if result != tc.expected {
			t.Errorf("AuthenticationContext.HasClaimValue(%s, %s) failed: got %v, expected %v", tc.name, tc.value, result, tc.expected)
		}
	}
}

func Test_AuthenticationContext_HasRole(t *testing.T) {
	type testCase struct {
		value    string
		expected bool
	}
	tests := []testCase{
		{"role1", true},
		{"ROLE1", true},
		{"role3", false},
	}
	ctx := &authenticationContext{
		authenticated: true,
		claims: ClaimMap{ClaimRole: &Claim{
			ClaimRole, []string{"role1", "role2"},
		}},
	}

	for _, tc := range tests {
		result := ctx.HasRole(tc.value)
		if result != tc.expected {
			t.Errorf("AuthenticationContext.HasRole(%s) failed: got %v, expected :%v", tc.value, result, tc.expected)
		}
	}
}

func Test_AuthenticationContext_HasScope(t *testing.T) {
	type testCase struct {
		value    string
		expected bool
	}
	tests := []testCase{
		{"scope1", true},
		{"SCOPE2", true},
		{"scope3", false},
	}
	ctx := &authenticationContext{
		authenticated: true,
		claims: ClaimMap{ClaimScope: &Claim{
			ClaimScope, []string{"scope1", "scope2"},
		}},
	}

	for _, tc := range tests {
		result := ctx.HasScope(tc.value)
		if result != tc.expected {
			t.Errorf("AuthenticationContext.HasScope(%s) failed: got %v, expected :%v", tc.value, result, tc.expected)
		}
	}
}
