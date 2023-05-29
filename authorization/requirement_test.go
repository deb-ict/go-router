package authorization

import (
	"testing"

	"github.com/deb-ict/go-router/authentication"
)

type combinedRequirementMock struct {
	result bool
}

func (r *combinedRequirementMock) MeetsRequirement(auth authentication.Context) bool {
	return r.result
}

func testUserContext() authentication.Context {
	claims := make(authentication.ClaimMap)
	claims.SetSubjectId("1")
	claims.SetName("test_user")
	claims.AddRoles("member")
	claims.AddScopes("api.read", "api.write")
	claims.AddClaim("team", "teamx")
	return authentication.NewContext(true, claims)
}

func testEmptyContext() authentication.Context {
	claims := make(authentication.ClaimMap)
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
		{"empty", testEmptyContext(), false},
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
	tests := []testCase{
		{"anonymous", authentication.AnonymouseContext(), authentication.ClaimName, []string{"anonymous"}, false},
		{"name", testUserContext(), authentication.ClaimName, []string{"test_user"}, true},
		{"team", testUserContext(), "team", []string{"teamx"}, true},
	}

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
		auth     authentication.Context
		roles    []string
		expected bool
	}
	tests := []testCase{
		{"anonymous", authentication.AnonymouseContext(), []string{"member", "admin"}, false},
		{"with_role", testUserContext(), []string{"member", "admin"}, true},
		{"without_role", testUserContext(), []string{"role1", "role2"}, false},
	}

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
		auth     authentication.Context
		scopes   []string
		expected bool
	}
	tests := []testCase{
		{"anonymous", authentication.AnonymouseContext(), []string{"api.read", "api.delete"}, false},
		{"with_scope", testUserContext(), []string{"api.read", "api.delete"}, true},
		{"without_scope", testUserContext(), []string{"api.test", "api.delete"}, false},
	}

	for _, tc := range tests {
		req := NewScopeRequirement(tc.scopes...)
		result := req.MeetsRequirement(tc.auth)
		if tc.expected != result {
			t.Errorf("UserRequirement(%s) failed: got %v, expected %v", tc.name, result, tc.expected)
		}
	}
}

func Test_CombinedRequirement(t *testing.T) {
	positive := &combinedRequirementMock{result: true}
	negative := &combinedRequirementMock{result: false}

	type testCase struct {
		name         string
		auth         authentication.Context
		requireAll   bool
		requirements []Requirement
		expected     bool
	}
	tests := []testCase{
		{"anonymous", authentication.AnonymouseContext(), false, []Requirement{}, false},
		{"nil_requirements", testUserContext(), false, nil, true},
		{"require_any_positive_negative", testUserContext(), false, []Requirement{positive, negative}, true},
		{"require_any_negative_positive", testUserContext(), false, []Requirement{negative, positive}, true},
		{"require_any_negative_negative", testUserContext(), false, []Requirement{negative, negative}, false},
		{"require_all_positive_negative", testUserContext(), true, []Requirement{positive, negative}, false},
		{"require_all_negative_positive", testUserContext(), true, []Requirement{negative, positive}, false},
		{"require_all_positive_positive", testUserContext(), true, []Requirement{positive, positive}, true},
	}

	for _, tc := range tests {
		req := NewCombinedRequirement(tc.requireAll, tc.requirements...)
		result := req.MeetsRequirement(tc.auth)
		if tc.expected != result {
			t.Errorf("UserRequirement(%s) failed: got %v, expected %v", tc.name, result, tc.expected)
		}
	}
}
