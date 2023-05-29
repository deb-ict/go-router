package authorization

import (
	"testing"

	"github.com/deb-ict/go-router/authentication"
)

type policyRequirementMock struct {
	calls  *int
	result bool
}

func (r *policyRequirementMock) MeetsRequirement(auth authentication.Context) bool {
	*r.calls++
	return r.result
}

func Test_NewPolicy(t *testing.T) {
	requirement := &policyRequirementMock{}
	policy := NewPolicy("test", requirement)
	if policy.GetName() != "test" {
		t.Error("NewPolicy() failed: Name not set")
	}
	if len(policy.GetRequirements()) == 1 {
		if policy.GetRequirements()[0] != requirement {
			t.Error("NewPolicy() failed: Requirement not found")
		}
	} else {
		t.Error("NewPolicy() failed: Requirements not set")
	}
}

func Test_Policy_GetRequirements_InitializeSlice(t *testing.T) {
	policy := &policy{
		requirements: nil,
	}
	requirements := policy.GetRequirements()
	if requirements == nil {
		t.Error("Policy.GetRequirements() failed: nil requirements scile not initialized")
	}
}

func Test_Policy_MeetsRequirements(t *testing.T) {
	context := authentication.AnonymouseContext()

	positive := &policyRequirementMock{
		result: true,
	}
	negative := &policyRequirementMock{
		result: false,
	}

	type testCase struct {
		policy         *policy
		context        authentication.Context
		calls          int
		expectedResult bool
		expectedCalls  int
	}
	tests := []testCase{
		{&policy{name: "nil context", requirements: []Requirement{}}, nil, 0, false, 0},
		{&policy{name: "no requirements", requirements: []Requirement{}}, context, 0, true, 0},
		{&policy{name: "positive", requirements: []Requirement{positive}}, context, 0, true, 1},
		{&policy{name: "negative", requirements: []Requirement{negative}}, context, 0, false, 1},
		{&policy{name: "positive_positive", requirements: []Requirement{positive, positive}}, context, 0, true, 2},
		{&policy{name: "positive_negative", requirements: []Requirement{positive, negative}}, context, 0, false, 2},
		{&policy{name: "negative_negative", requirements: []Requirement{negative, negative}}, context, 0, false, 1},
		{&policy{name: "negative_positive", requirements: []Requirement{negative, positive}}, context, 0, false, 1},
	}

	for _, tc := range tests {
		for _, r := range tc.policy.requirements {
			mock := r.(*policyRequirementMock)
			mock.calls = &tc.calls
		}

		result := tc.policy.MeetsRequirements(tc.context)
		if result != tc.expectedResult {
			t.Errorf("Policy.MeetsRequirements(%s) failed: invalid result: got %v, expected %v", tc.policy.name, result, tc.expectedResult)
		}
		if tc.calls != tc.expectedCalls {
			t.Errorf("Policy.MeetsRequirements(%s) failed: invalid number of checks: got %v, expected %v", tc.policy.name, tc.calls, tc.expectedCalls)
		}
	}
}
