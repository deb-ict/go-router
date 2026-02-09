package router

import (
	"net/http"
	"strings"
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
		{"/api/v1/{id}", false, false, []string{"api", "v1", "id"}},
		{"api/v1", false, true, []string{}},
		{"", false, true, []string{}},
		{"api/v1", false, true, []string{}},
	}

	for _, tc := range tests {
		root := &Node{}
		node := root.BuildTree(tc.pattern)
		if tc.isNil && node != nil {
			t.Errorf("Node.BuildTree(%s) failed: Result node should be <nil>", tc.pattern)
		}
		if !tc.isNil && node == nil {
			t.Errorf("Node.BuildTree(%s) failed: Result node should not be <nil>", tc.pattern)
		}
		if tc.isRoot && node != root {
			t.Errorf("Node.BuildTree(%s) failed: Result node should be root node", tc.pattern)
		}
		if len(tc.segments) > 0 {
			exp := tc.segments[len(tc.segments)-1]
			if node == nil || node.Segment != exp {
				t.Errorf("Node.BuildTree(%s) failed: Invalid result node: got %s, expected %s", tc.pattern, nodeName(node), exp)
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
		pattern    string
		isRoot     bool
		isNil      bool
		expected   string
		paramKey   string
		paramValue string
		segments   []string
	}
	tests := []testCase{
		{"/", true, false, "", "", "", []string{""}},
		{"//", true, false, "", "", "", []string{""}},
		{"/api/v1", false, false, "v1", "", "", []string{"api", "v1"}},
		{"/api/v1", false, false, "v1", "", "", []string{"api", "v1", "test"}},
		{"/api//v1", false, false, "v1", "", "", []string{"api", "v1"}},
		{"/api/v1/", false, false, "v1", "", "", []string{"api", "v1"}},
		{"/api/v1/123", false, false, "id", "id", "123", []string{"api", "v1", "{id}"}},
		{"api/v1", false, true, "", "", "", []string{"api", "v1"}},
		{"", false, true, "", "", "", []string{""}},
		{"api/v1", false, true, "", "", "", []string{""}},
	}

	for _, tc := range tests {
		params := make(map[string]string)
		root := buildTestTree(tc.segments)
		node := root.FindNode(tc.pattern, params)
		if tc.isNil && node != nil {
			t.Errorf("Node.FindNode(%s) failed: Result node should be <nil>", tc.pattern)
		}
		if !tc.isNil && node == nil {
			t.Errorf("Node.FindNode(%s) failed: Result node should not be <nil>", tc.pattern)
		}
		if tc.isRoot && node != root {
			t.Errorf("Node.FindNode(%s) failed: Result node should be root node", tc.pattern)
		}
		if node != nil && node.Segment != tc.expected {
			t.Errorf("Node.FindNode(%s) failed: Invalid result node: got %s, expected %s", tc.pattern, nodeName(node), tc.expected)
		}
		if tc.paramKey != "" {
			param := params[tc.paramKey]
			if param != tc.paramValue {
				t.Errorf("Node.FindNode(%s) failed: Param not found: got %s, expected %s", tc.pattern, param, tc.paramValue)
			}
		}
	}
}

func Test_Node_getPath(t *testing.T) {
	type testCase struct {
		pattern  string
		expected string
	}
	tests := []testCase{
		{"/api/v1", "/api/v1"},
		{"/api/v1?query", "/api/v1"},
	}

	for _, tc := range tests {
		node := &Node{}
		result := node.getPath(tc.pattern)
		if result != tc.expected {
			t.Errorf("Node.getPath(%s) failed: got %v, expected %v", tc.pattern, result, tc.expected)
		}
	}
}

func Test_Node_buildSegment(t *testing.T) {
	type testCase struct {
		node             *Node
		segments         []string
		expectedExisting bool
		expectedNode     *Node
	}
	root := &Node{Segment: "self", Type: NodeTypePath}
	node := &Node{Segment: "next", Type: NodeTypePath}
	nodeWithHandler := &Node{
		Segment: "next",
		Routes: []*Route{
			{Handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})},
		},
	}
	rootWithChild := &Node{
		Segment: "self",
		Nodes:   []*Node{node},
	}
	rootWithHandlerChild := &Node{
		Segment: "self",
		Nodes:   []*Node{nodeWithHandler},
	}
	tests := []testCase{
		{root, []string{}, true, root},
		{root, []string{""}, true, root},
		{root, []string{"", "next"}, false, node},
		{rootWithChild, []string{"next"}, true, node},
		{rootWithHandlerChild, []string{"next"}, false, node},
	}

	for _, tc := range tests {
		result := tc.node.buildSegment(tc.segments)
		pattern := strings.Join(tc.segments, "/")
		if tc.expectedExisting {
			if tc.expectedExisting && result != tc.expectedNode {
				t.Errorf("Node.buildSegment(%s) failed: get %v, expected existing %v", pattern, nodeName(result), nodeName(tc.expectedNode))
			}
		} else {
			if result == nil {
				t.Errorf("Node.buildSegment(%s) failed: got <nil>, expected new %v", pattern, nodeName(tc.expectedNode))
			} else if result.Segment != tc.expectedNode.Segment {
				t.Errorf("Node.buildSegment(%s) failed: got %v, expected new %v", pattern, nodeName(result), nodeName(tc.expectedNode))
			}
		}
	}
}

func Test_Node_findSegment(t *testing.T) {
	type testCase struct {
		node     *Node
		segments []string
		expected *Node
	}
	node := &Node{Segment: "next", Type: NodeTypePath}
	root := &Node{Segment: "self", Nodes: []*Node{node}}
	tests := []testCase{
		{root, []string{}, root},
		{root, []string{""}, root},
		{root, []string{"next"}, node},
		{root, []string{"", "next"}, node},
		{root, []string{"notfound"}, nil},
	}

	for _, tc := range tests {
		params := make(map[string]string)
		result := tc.node.findSegment(tc.segments, params)
		pattern := strings.Join(tc.segments, "/")
		if result != tc.expected {
			t.Errorf("Node.findSegment(%s) failed: got %v, expected: %v", pattern, nodeName(result), nodeName(tc.expected))
		}
	}
}

func Test_Node_findChildSegment(t *testing.T) {
	type testCase struct {
		node     *Node
		segment  []string
		expected *Node
	}
	pathNode := &Node{Segment: "path", Type: NodeTypePath}
	paramNode := &Node{Segment: "{param}", Type: NodeTypeParam}
	tests := []testCase{
		{&Node{Nodes: []*Node{pathNode}}, []string{"path"}, pathNode},
		{&Node{Nodes: []*Node{paramNode}}, []string{"path"}, paramNode},
		{&Node{Nodes: []*Node{pathNode}}, []string{"notfound"}, nil},
	}

	for _, tc := range tests {
		params := make(map[string]string)
		result := tc.node.findChildSegment(tc.segment[0], tc.segment, params)
		if result != tc.expected {
			t.Errorf("Node.findChildSegment(%s) failed: got %v, expected %v", tc.segment, nodeName(result), nodeName(tc.expected))
		}
	}
}

func Test_Node_getNodeType(t *testing.T) {
	type testCase struct {
		segment  string
		expected NodeType
	}
	tests := []testCase{}

	for _, tc := range tests {
		node := &Node{}
		result := node.getNodeType(tc.segment)
		if result != tc.expected {
			t.Errorf("Node.getNodeType(%s) failed: got %v, expected %v", tc.segment, result, tc.expected)
		}
	}
}

func Test_Node_getNodeValue(t *testing.T) {
	type testCase struct {
		segment  string
		nodeType NodeType
		expected string
	}
	tests := []testCase{
		{"api", NodeTypePath, "api"},
		{"{id}", NodeTypeParam, "id"},
	}

	for _, tc := range tests {
		node := &Node{}
		result := node.getNodeValue(tc.segment, tc.nodeType)
		if result != tc.expected {
			t.Errorf("Node.getNodeValue(%s) failed: got %s, expected %s", tc.segment, result, tc.expected)
		}
	}
}

func Test_Node_isParamSegment(t *testing.T) {
	type testCase struct {
		segment  string
		expected bool
	}
	tests := []testCase{
		{"api", false},
		{"{param}", true},
		{"a{param}b", false},
	}

	for _, tc := range tests {
		node := &Node{}
		result := node.isParamSegment(tc.segment)
		if result != tc.expected {
			t.Errorf("Node.isParamSegment(%s) failed: got %v, expected %v", tc.segment, result, tc.expected)
		}
	}
}

func Test_NodeType_String(t *testing.T) {
	path := NodeTypePath.String()
	if path != "path" {
		t.Errorf("NodeTypePath.String() failed: got '%v', expected 'path'", path)
	}

	param := NodeTypeParam.String()
	if param != "param" {
		t.Errorf("NodeTypeParam.String() failed: got '%v', expected 'param'", param)
	}

	undefined := NodeTypeUndefined.String()
	if undefined != "undefined" {
		t.Errorf("NodeTypeUndefined.String() failed: got '%v', expected 'undefined'", undefined)
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
	node := &Node{}
	node.Type = node.getNodeType(segment)
	node.Segment = node.getNodeValue(segment, node.Type)
	parent.Nodes = []*Node{
		node,
	}
	return node
}

func validateDepth(t *testing.T, root *Node, pattern string, depth int) {
	node := root
	for i := 0; i < depth; i++ {
		if node.Nodes == nil || len(node.Nodes) == 0 {
			t.Errorf("Node.BuildTree(%s) failed: Node %s has no child nodes, expected at least 1", pattern, node.Segment)
			return
		}
		if node.Nodes[0].Parent != node {
			t.Errorf("Node.BuildTree(%s) failed: Node parent incorrect", pattern)
		}
		node = node.Nodes[0]
	}

	if node.Nodes != nil && len(node.Nodes) > 0 {
		t.Errorf("Node.BuildTree(%s) failed: Node %s has child nodes, expected 0", pattern, node.Segment)
	}
}

func validateChildNode(t *testing.T, root *Node, pattern string, depth int, segment string) {
	node := root
	for i := 0; i < depth; i++ {
		if node.Nodes == nil || len(node.Nodes) == 0 {
			t.Errorf("Node.BuildTree(%s) failed: Node %s has no child nodes, expected at least 1", pattern, node.Segment)
			return
		}
		node = node.Nodes[0]
	}
	if node.Segment != segment {
		t.Errorf("Node.BuildTree(%s) failed: Segment not found: got %s, expected %s", pattern, node.Segment, segment)
	}
}

func nodeName(n *Node) string {
	if n == nil {
		return "<nil>"
	}
	if n.Segment == "" {
		return "<empty>"
	}
	return n.Segment
}
