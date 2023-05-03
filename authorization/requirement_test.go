package authorization

import (
	"testing"

	"github.com/deb-ict/go-router/authentication"
)

func testUserContext() authentication.Context {
	claims := make(authentication.ClaimMap)
	claims.SetSubjectId("1")
	claims.SetName("test_user")
	claims.AddRoles("member")
	claims.AddScopes("api.read", "api.write")
	claims.AddClaim("team", "teamx")
	return authentication.NewContext(true, claims)
}

func Test_UserRequirement(t *testing.T) {
	type testCase struct {
		name     string
		auth     authentication.Context
		expected bool
	}
	tests := []testCase{
		{"anonymous", authentication.AnonymouseContext(), false},
		{"user", testUserContext(), true},
	}

	for _, tc := range tests {
		req := NewUserRequirement()
		result := req.MeetsRequirement(tc.auth)
		if tc.expected != result {
			t.Errorf("UserRequirement(%s) failed: got %v, expected %v", tc.name, result, tc.expected)
		}
	}
}

func Test_ClaimRequirement(t *testing.T) {
	type testCase struct {
		name     string
		auth     authentication.Context
		claim    string
		values   []string
		expected bool
	}
	tests := []testCase{}

	for _, tc := range tests {
		req := NewClaimRequirement(tc.claim, tc.values...)
		result := req.MeetsRequirement(tc.auth)
		if tc.expected != result {
			t.Errorf("UserRequirement(%s) failed: got %v, expected %v", tc.name, result, tc.expected)
		}
	}
}

func Test_RoleRequirement(t *testing.T) {
	type testCase struct {
		name     string
		roles    []string
		auth     authentication.Context
		expected bool
	}
	tests := []testCase{}

	for _, tc := range tests {
		req := NewRoleRequirement(tc.roles...)
		result := req.MeetsRequirement(tc.auth)
		if tc.expected != result {
			t.Errorf("UserRequirement(%s) failed: got %v, expected %v", tc.name, result, tc.expected)
		}
	}
}

func Test_ScopeRequirement(t *testing.T) {
	type testCase struct {
		name     string
		scopes   []string
		auth     authentication.Context
		expected bool
	}
	tests := []testCase{}

	for _, tc := range tests {
		req := NewScopeRequirement(tc.scopes...)
		result := req.MeetsRequirement(tc.auth)
		if tc.expected != result {
			t.Errorf("UserRequirement(%s) failed: got %v, expected %v", tc.name, result, tc.expected)
		}
	}
}

func Test_CombinedRequirement(t *testing.T) {
	type testCase struct {
		name     string
		auth     authentication.Context
		expected bool
	}
	tests := []testCase{}

	for _, tc := range tests {
		req := NewUserRequirement()
		result := req.MeetsRequirement(tc.auth)
		if tc.expected != result {
			t.Errorf("UserRequirement(%s) failed: got %v, expected %v", tc.name, result, tc.expected)
		}
	}
}
