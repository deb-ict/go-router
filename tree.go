package router

import (
	"strings"
)

type NodeType int

const (
	NodeTypeUndefined NodeType = iota
	NodeTypePath
	NodeTypeParam
)

type Node struct {
	Segment string
	Type    NodeType
	Route   *Route
	Nodes   []*Node
}

func (n *Node) BuildTree(pattern string) *Node {
	if pattern == "" {
		return nil
	}
	if pattern[0] != '/' {
		return nil
	}

	segments := strings.Split(n.getPath(pattern), "/")
	if len(segments) == 0 {
		return nil
	}
	return n.buildSegment(segments[1:])
}

func (n *Node) FindNode(pattern string) *Node {
	if pattern == "" {
		return nil
	}
	if pattern[0] != '/' {
		return nil
	}

	segments := strings.Split(n.getPath(pattern), "/")
	if len(segments) == 0 {
		return nil
	}
	return n.findSegment(segments[1:])
}

func (n *Node) getPath(pattern string) string {
	i := strings.Index(pattern, "?")
	if i == -1 {
		return pattern
	}
	return pattern[:i]
}

func (n *Node) buildSegment(segments []string) *Node {
	// Validate the segments
	numSegments := len(segments)
	if numSegments == 0 {
		return n
	}
	if segments[0] == "" && numSegments == 1 {
		return n
	}
	if segments[0] == "" && numSegments > 1 {
		return n.buildSegment(segments[1:])
	}
	segment := strings.ToLower(segments[0])

	// Find existing node
	var node *Node
	for _, cn := range n.Nodes {
		if cn.Segment == segment {
			node = cn
			break
		}
	}

	// Create a node for the segment and append to parent
	if node == nil || (node.hasHandler() && numSegments == 1) {
		node = &Node{
			Segment: segment,
			Type:    NodeTypePath,
		}
		//TODO: Check if segment is parameter: {name} and change type to NodeTypeParam
		if n.Nodes == nil {
			n.Nodes = make([]*Node, 0)
		}
		n.Nodes = append(n.Nodes, node)
	}

	// Build recursive nodes
	return node.buildSegment(segments[1:])
}

func (n *Node) findSegment(segments []string) *Node {
	// Validate the segments
	numSegments := len(segments)
	if numSegments == 0 {
		return n
	}
	if segments[0] == "" && numSegments == 1 {
		return n
	}
	if segments[0] == "" && numSegments > 1 {
		return n.findSegment(segments[1:])
	}

	if n.Nodes != nil {
		for _, node := range n.Nodes {
			if node.Segment == strings.ToLower(segments[0]) {
				return node.findSegment(segments[1:])
			}
		}
	}
	return nil
}

func (n *Node) hasRoute() bool {
	return n.Route != nil
}

func (n *Node) hasHandler() bool {
	return n.hasRoute() && n.Route.Handler != nil
}

func (t NodeType) String() string {
	switch t {
	case NodeTypePath:
		return "path"
	case NodeTypeParam:
		return "param"
	}
	return "undefined"
}
