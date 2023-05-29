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
		expected *Claim
		added    bool
	}
	tests := []testCase{
		{ClaimMap{"test1": &Claim{"test1", []string{"value1", "value2"}}}, &Claim{"test1", []string{"value1", "value2"}}, false},
		{ClaimMap{"test1": &Claim{"test1", []string{"value1", "value2"}}}, &Claim{"test2", []string{}}, true},
	}

	for _, tc := range tests {
		result := tc.claims.GetClaim(tc.expected.Name)
		compareClaims(t, "ClaimMap.GetClaim()", result, tc.expected)
		if tc.added && len(tc.claims) == 1 {
			t.Errorf("ClaimMap.GetClaim() failed: returned claim not added to map")
		}
		if !tc.added && len(tc.claims) > 1 {
			t.Errorf("ClaimMap.GetClaim() failed: returned claim added to map")
		}
	}
}

func Test_ClaimMap_AddClaim(t *testing.T) {
	type testCase struct {
		claims   ClaimMap
		name     string
		value    string
		expected []string
	}
	tests := []testCase{
		{ClaimMap{"test1": &Claim{"test1", []string{"value1", "value2"}}}, "test1", "value3", []string{"value1", "value2", "value3"}},
		{ClaimMap{"test1": &Claim{"test1", []string{"value1", "value2"}}}, "test1", "value1", []string{"value1", "value2"}},
		{ClaimMap{"test1": &Claim{"test1", []string{"value1", "value2"}}}, "test2", "value1", []string{"value1"}},
	}

	for _, tc := range tests {
		tc.claims.AddClaim(tc.name, tc.value)
		result, ok := tc.claims[tc.name]
		if !ok {
			t.Errorf("ClaimMap.AddClaim() failed: Claim not found in map")
		}
		compareClaimValues(t, "ClaimMap.AddClaim()", result.Values, tc.expected)
	}
}

func Test_ClaimMap_SetClaimSingleValue(t *testing.T) {
	type testCase struct {
		claims   ClaimMap
		name     string
		value    string
		expected []string
	}
	tests := []testCase{
		{ClaimMap{"test1": &Claim{"test1", []string{"value1", "value2"}}}, "test1", "newvalue", []string{"newvalue"}},
		{ClaimMap{"test1": &Claim{"test1", []string{"value1", "value2"}}}, "test2", "newvalue", []string{"newvalue"}},
	}

	for _, tc := range tests {
		tc.claims.SetClaimSingleValue(tc.name, tc.value)
		result, ok := tc.claims[tc.name]
		if !ok {
			t.Errorf("ClaimMap.SetClaimSingleValue() failed: Claim not found in map")
		}
		compareClaimValues(t, "ClaimMap.SetClaimSingleValue()", result.Values, tc.expected)
	}
}

func Test_ClaimMap_AddRoles(t *testing.T) {
	claims := ClaimMap{}
	claims.AddRoles("role1", "role2", "role1")

	claim, ok := claims[ClaimRole]
	if !ok {
		t.Fatal("ClaimMap.AddRoles() failed: Claim not found")
	}
	if len(claim.Values) != 2 {
		t.Fatalf("ClaimMap.AddRoles() failed: Incorrect number of values: got %v, expected 2", len(claim.Values))
	}
	if claim.Values[0] != "role1" {
		t.Errorf("ClaimMap.AddRoles() failed: Incorrect role name at index 0: got %v, expected role1", claim.Values[0])
	}
	if claim.Values[1] != "role2" {
		t.Errorf("ClaimMap.AddRoles() failed: Incorrect role name at index 1: got %v, expected role2", claim.Values[1])
	}
}

func Test_ClaimMap_AddScopes(t *testing.T) {
	claims := ClaimMap{}
	claims.AddScopes("scope1", "scope2", "scope1")

	claim, ok := claims[ClaimScope]
	if !ok {
		t.Fatal("ClaimMap.AddScopes() failed: Claim not found")
	}
	if len(claim.Values) != 2 {
		t.Fatalf("ClaimMap.AddScopes() failed: Incorrect number of values: got %v, expected 2", len(claim.Values))
	}
	if claim.Values[0] != "scope1" {
		t.Errorf("ClaimMap.AddScopes() failed: Incorrect scope name at index 0: got %v, expected scope1", claim.Values[0])
	}
	if claim.Values[1] != "scope2" {
		t.Errorf("ClaimMap.AddScopes() failed: Incorrect scope name at index 1: got %v, expected scope2", claim.Values[1])
	}
}

func Test_ClaimMap_SetSubjectId(t *testing.T) {
	claims := ClaimMap{}
	claims.SetSubjectId("123")

	claim, ok := claims[ClaimSubjectId]
	if !ok {
		t.Fatal("ClaimMap.SetSubjectId() failed: Claim not found")
	}
	if len(claim.Values) != 1 {
		t.Fatalf("ClaimMap.SetSubjectId() failed: Incorrect number of values: got %v, expected 1", len(claim.Values))
	}
	if claim.Values[0] != "123" {
		t.Errorf("ClaimMap.SetSubjectId() failed: Incorrect value: got %v, expected 123", claim.Values[0])
	}
}

func Test_ClaimMap_SetName(t *testing.T) {
	claims := ClaimMap{}
	claims.SetName("tester")

	claim, ok := claims[ClaimName]
	if !ok {
		t.Fatal("ClaimMap.SetName() failed: Claim not found")
	}
	if len(claim.Values) != 1 {
		t.Fatalf("ClaimMap.SetName() failed: Incorrect number of values: got %v, expected 1", len(claim.Values))
	}
	if claim.Values[0] != "tester" {
		t.Errorf("ClaimMap.SetName() failed: Incorrect value: got %v, expected tester", claim.Values[0])
	}
}

func compareClaims(t *testing.T, testName string, result *Claim, expected *Claim) {
	if result.Name != expected.Name {
		t.Errorf("%s failed: name not equal: got %v, expected %v", testName, result.Name, expected.Name)
	}
	compareClaimValues(t, testName, result.Values, expected.Values)
}

func compareClaimValues(t *testing.T, testName string, result []string, expected []string) {
	if len(result) != len(expected) {
		t.Errorf("%s failed: number of values not equal: got %v, expected %v", testName, len(result), len(expected))
	} else {
		for i := 0; i < len(result); i++ {
			if result[i] != expected[i] {
				t.Errorf("%s failed: values[%d] not equal: got %v, expected %v", testName, i, result[i], expected[i])
			}
		}
	}
}
