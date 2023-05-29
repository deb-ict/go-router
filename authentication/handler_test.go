package authentication

import (
	"testing"
)

func Test_equalFold(t *testing.T) {
	type testCase struct {
		s string
		t string
		e bool
	}
	tests := []testCase{
		{"abc", "x", false},
		{"abc", "xyz", false},
		{"abc", "abc", true},
		{"abc", "ABC", true},
		{"ABC", "abc", true},
	}

	for _, tc := range tests {
		result := equalFold(tc.s, tc.t)
		if result != tc.e {
			t.Errorf("equalFold(%s, %s) failed: got %v, expected %v", tc.s, tc.t, result, tc.e)
		}
	}
}
