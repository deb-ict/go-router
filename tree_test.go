package router

import (
	"testing"
)

func Test_Node_BuildTree(t *testing.T) {
	type testCase struct {
		pattern  string
		isRoot   bool
		isNil    bool
		segments []string
	}
	tests := []testCase{
		{"/", true, false, []string{}},
		{"//", true, false, []string{}},
		{"/api/v1", false, false, []string{"api", "v1"}},
		{"/api/v1?a=b", false, false, []string{"api", "v1"}},
		{"/api//v1", false, false, []string{"api", "v1"}},
		{"/api/v1/", false, false, []string{"api", "v1"}},
		{"/api/v1/{id}", false, false, []string{"api", "v1", "{id}"}},
		{"api/v1", false, true, []string{}},
	}

	for _, tc := range tests {
		root := &Node{}
		node := root.BuildTree(tc.pattern)
		if tc.isNil && node != nil {
			t.Errorf("[%s] Result node should be nil", tc.pattern)
		}
		if !tc.isNil && node == nil {
			t.Errorf("[%s] Result node should not be nil", tc.pattern)
		}
		if tc.isRoot && node != root {
			t.Errorf("[%s] Result node should be root node", tc.pattern)
		}
		if len(tc.segments) > 0 {
			exp := tc.segments[len(tc.segments)-1]
			if node == nil || node.Segment != exp {
				t.Errorf("[%s] Invalid result node: got %s, expected %s", tc.pattern, node.Segment, exp)
			}
		}
		validateDepth(t, root, tc.pattern, len(tc.segments))
		for i, exp := range tc.segments {
			validateChildNode(t, root, tc.pattern, i+1, exp)
		}
	}
}

func Test_Node_FindNode(t *testing.T) {
	type testCase struct {
		pattern  string
		isRoot   bool
		isNil    bool
		expected string
		segments []string
	}
	tests := []testCase{
		{"/", true, false, "", []string{""}},
		{"//", true, false, "", []string{""}},
		{"/api/v1", false, false, "v1", []string{"api", "v1"}},
		{"/api/v1", false, false, "v1", []string{"api", "v1", "test"}},
		{"/api//v1", false, false, "v1", []string{"api", "v1"}},
		{"/api/v1/", false, false, "v1", []string{"api", "v1"}},
		{"/api/v1/{id}", false, false, "{id}", []string{"api", "v1", "{id}"}},
		{"api/v1", false, true, "", []string{"api", "v1"}},
	}

	for _, tc := range tests {
		root := buildTestTree(tc.segments)
		node := root.FindNode(tc.pattern)
		if tc.isNil && node != nil {
			t.Errorf("[%s] Result node should be nil", tc.pattern)
		}
		if !tc.isNil && node == nil {
			t.Errorf("[%s] Result node should not be nil", tc.pattern)
		}
		if tc.isRoot && node != root {
			t.Errorf("[%s] Result node should be root node", tc.pattern)
		}
		if node != nil && node.Segment != tc.expected {
			t.Errorf("[%s] Invalid result node: got %s, expected %s", tc.pattern, node.Segment, tc.expected)
		}
	}
}

func buildTestTree(segments []string) *Node {
	root := &Node{}
	node := root
	for _, s := range segments {
		node = appendTestNode(node, s)
	}
	return root
}

func appendTestNode(parent *Node, segment string) *Node {
	node := &Node{
		Segment: segment,
	}
	parent.Nodes = []*Node{
		node,
	}
	return node
}

func validateDepth(t *testing.T, root *Node, pattern string, depth int) {
	node := root
	for i := 0; i < depth; i++ {
		if node.Nodes == nil || len(node.Nodes) == 0 {
			t.Errorf("[%s] Node %s has no child nodes, expected at least 1", pattern, node.Segment)
			return
		}
		node = node.Nodes[0]
	}

	if node.Nodes != nil && len(node.Nodes) > 0 {
		t.Errorf("[%s] Node %s has child nodes, expected 0", pattern, node.Segment)
	}
}

func validateChildNode(t *testing.T, root *Node, pattern string, depth int, segment string) {
	node := root
	for i := 0; i < depth; i++ {
		if node.Nodes == nil || len(node.Nodes) == 0 {
			t.Errorf("[%s] Node %s has no child nodes, expected at least 1", pattern, node.Segment)
			return
		}
		node = node.Nodes[0]
	}
	if node.Segment != segment {
		t.Errorf("[%s] Segment not found: got %s, expected %s", pattern, node.Segment, segment)
	}
}
