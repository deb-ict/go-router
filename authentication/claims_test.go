package authentication

import (
	"testing"
)

func Test_Claim_First(t *testing.T) {
	type testCase struct {
		claim    Claim
		expected string
	}
	tests := []testCase{
		{Claim{Values: []string{"value1", "value2"}}, "value1"},
		{Claim{Values: []string{}}, ""},
		{Claim{Values: nil}, ""},
	}

	for _, tc := range tests {
		var result = tc.claim.First()
		if result != tc.expected {
			t.Errorf("Claim.First() failed: got %v, expected %v", result, tc.expected)
		}
	}
}

func Test_Claim_Value(t *testing.T) {
	type testCase struct {
		claim    Claim
		index    int
		expected string
	}
	tests := []testCase{
		{Claim{Values: []string{"value1", "value2"}}, 0, "value1"},
		{Claim{Values: []string{"value1", "value2"}}, 1, "value2"},
		{Claim{Values: []string{"value1", "value2"}}, 2, ""},
		{Claim{Values: []string{}}, 0, ""},
		{Claim{Values: nil}, 0, ""},
	}

	for _, tc := range tests {
		var result = tc.claim.Value(tc.index)
		if result != tc.expected {
			t.Errorf("Claim.Value(%d) failed: got %v, expected %v", tc.index, result, tc.expected)
		}
	}
}

func Test_Claim_HasValue(t *testing.T) {
	type testCase struct {
		claim    Claim
		value    string
		expected bool
	}
	tests := []testCase{
		{Claim{Values: []string{"value1", "value2"}}, "value1", true},
		{Claim{Values: []string{"value1", "value2"}}, "value2", true},
		{Claim{Values: []string{"value1", "value2"}}, "value3", false},
		{Claim{Values: []string{"value1", "value2"}}, "VALUE1", true},
		{Claim{Values: []string{"value1", "value2"}}, "", false},
		{Claim{Values: []string{}}, "value1", false},
		{Claim{Values: nil}, "value1", false},
	}

	for _, tc := range tests {
		var result = tc.claim.HasValue(tc.value)
		if result != tc.expected {
			t.Errorf("Claim.Value(%s) failed: got %v, expected %v", tc.value, result, tc.expected)
		}
	}
}

func Test_ClaimMap_GetClaim(t *testing.T) {
	type testCase struct {
		claims   ClaimMap
		expected Claim
		added    bool
	}
	claim1 := Claim{"test1", []string{"value1", "value2"}}
	claim2 := Claim{"test2", []string{"value1", "value2"}}
	tests := []testCase{
		{ClaimMap{claim1.Name: claim1}, claim1, false},
		{ClaimMap{claim1.Name: claim1}, claim2, true},
	}

	for _, tc := range tests {
		result := tc.claims.GetClaim(tc.expected.Name)
		if tc.added && len(tc.claims) == 1 {
			//not added
		}
	}
}

func Test_ClaimMap_AddClaim(t *testing.T) {

}

func Test_ClaimMap_SetClaimSingleValue(t *testing.T) {

}

func Test_ClaimMap_AddRoles(t *testing.T) {

}

func Test_ClaimMap_AddScopes(t *testing.T) {

}

func Test_ClaimMap_SetSubjectId(t *testing.T) {

}

func Test_ClaimMap_SetName(t *testing.T) {

}

func compareClaims(t *testing.T, a Claim, b Claim) bool {
	if a.Name != b.Name {
		return false
	}
	if len(a.Values) != len(b.Values) {
		return false
	}
	for i := 0; i < len(a.Values); i++ {
		if a.Values[i] != b.Values[i] {
			return false
		}
	}
	return true
}
