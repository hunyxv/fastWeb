package fastweb

import (
	"strings"
	"sync"
)

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

var nodeSlicePool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return make([]*node, 0)
	},
}

func releaseNodeSlice(nodes []*node) {
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

func (n *node) search(parts []string, depth int) *node {
	if len(parts) == depth || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[depth]
	children := n.matchChildren(part)
	defer releaseNodeSlice(children)

	for _, child := range children {
		result := child.search(parts, depth+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// -----------------------

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type nodeType uint8

const (
	static nodeType = iota // default
	root
	urlparam
	filepath
)

type newnode struct {
	pattern  string
	part     []rune
	children []*newnode
	isWild   bool
	nType    nodeType
	priority uint32
	handle   HandlerFunc
}

func (n *newnode) addRoute(pattern string, handle HandlerFunc) {
	n.priority++
	if len(n.pattern) == 0 && len(n.children) == 0 {
		n.nType = root
		n.insert(pattern, handle)
		return
	}

	var i int
	l := min(len(pattern), len(n.pattern))
	for j := 0; j < l; j++ {
		if n.pattern[j] == pattern[j] {
			i++
			continue
		}
		break
	}
	if i == len(pattern) {
		if n.handle != nil {
			panic("a handle is already registered for path '" + pattern + "'")
		}
	} else if i < len(n.pattern) {
		child := &newnode{
			pattern:  n.pattern[i:],
			isWild:   n.isWild,
			nType:    static,
			children: n.children,
			priority: n.priority - 1,
			handle: n.handle,
		}
		n.pattern = n.pattern[:i]
		n.children = []*newnode{child}

	}
}

func (n *newnode) insert(pattern string, handle HandlerFunc) {

}
