package router

import (
	"strings"
)

type NodeFilter func(n *Node)

type NodeType int

const (
	NodeTypeUndefined NodeType = iota
	NodeTypePath
	NodeTypeParam
)

type Node struct {
	Segment string
	Type    NodeType
	Parent  *Node
	Nodes   []*Node
	Routes  []*Route
}

func (n *Node) BuildTree(pattern string) *Node {
	if pattern == "" {
		return nil
	}
	if pattern[0] != '/' {
		return nil
	}

	segments := strings.Split(n.getPath(pattern), "/")
	return n.buildSegment(segments[1:])
}

func (n *Node) FindNode(pattern string, params RouteParams) *Node {
	if pattern == "" {
		return nil
	}
	if pattern[0] != '/' {
		return nil
	}

	segments := strings.Split(n.getPath(pattern), "/")
	return n.findSegment(segments[1:], params)
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
	nodeType := n.getNodeType(segment)
	nodeValue := n.getNodeValue(segment, nodeType)

	// Find existing node
	var node *Node
	for _, cn := range n.Nodes {
		if cn.Segment == nodeValue && cn.Type == nodeType {
			node = cn
			break
		}
	}

	// Create a node for the segment and append to parent
	if node == nil /*|| (node.hasHandler() && numSegments == 1)*/ {
		node = &Node{
			Segment: nodeValue,
			Type:    nodeType,
			Parent:  n,
		}
		if n.Nodes == nil {
			n.Nodes = make([]*Node, 0)
		}
		n.Nodes = append(n.Nodes, node)
	}

	// Build recursive nodes
	return node.buildSegment(segments[1:])
}

func (n *Node) findSegment(segments []string, params RouteParams) *Node {
	// Validate the segments
	numSegments := len(segments)
	if numSegments == 0 {
		return n
	}
	if segments[0] == "" && numSegments == 1 {
		return n
	}
	if segments[0] == "" && numSegments > 1 {
		return n.findSegment(segments[1:], params)
	}

	segment := strings.ToLower(segments[0])
	child := n.findChildSegment(segment, params)
	if child != nil {
		return child.findSegment(segments[1:], params)
	}
	return nil
}

func (n *Node) findChildSegment(segment string, params RouteParams) *Node {
	if n.Nodes != nil {
		for _, node := range n.Nodes {
			if node.Type == NodeTypeParam {
				params[node.Segment] = segment
				return node
			}
			if node.Type == NodeTypePath && node.Segment == segment {
				return node
			}
		}
	}
	return nil
}

func (n *Node) getNodeType(segment string) NodeType {
	if n.isParamSegment(segment) {
		return NodeTypeParam
	}
	return NodeTypePath
}

func (n *Node) getNodeValue(segment string, nodeType NodeType) string {
	if nodeType == NodeTypeParam {
		return segment[1 : len(segment)-1]
	}
	return segment
}

func (n *Node) isParamSegment(segment string) bool {
	return segment != "" && segment[0] == '{' && segment[len(segment)-1] == '}'
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
