package fastweb

import (
	"strings"
	"sync"
)

type node struct {
	pattern 	string
	part 		string
	children	[]*node
	isWild		bool
}


var nodeSlicePool *sync.Pool = &sync.Pool{
	New:	func() interface{} {
		return make([]*node, 0)
	},
}

func ReleaseNodeSlice(nodes []*node) {
	nodes = nodes[:0]
	nodeSlicePool.Put(nodes)
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := nodeSlicePool.Get().([]*node)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, depth int) {
	if len(parts) == depth {
		n.pattern = pattern
		return
	}

	part := parts[depth]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, depth+1)
}

func (n *node) search(parts[]string, depth int) *node {
	if len(parts) == depth || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[depth]
	children := n.matchChildren(part)
	defer ReleaseNodeSlice(children)

	for _, child := range children {
		result := child.search(parts, depth + 1)
		if result != nil {
			return result
		}
	}

	return nil
}