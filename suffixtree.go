// Ukkonen's algorithm to form suffixtree in linear time.
// Based mostly on http://marknelson.us/1996/08/01/suffix-trees/

package suffixtree

import (
	"strconv"
)

// A Node has some children (mapped by the first byte of the edge to that child),
// and an Edge (all the bytes for the edge to this Node from its parent).
type Node struct {
	Child  map[byte]*Node
	Edge   []byte
	suffix *Node
	name string
}

var nodeNum int

func newNode(edge []byte, child map[byte]*Node) *Node {
	name := strconv.Itoa(nodeNum)
	nodeNum++
	if child == nil {
		child = make(map[byte]*Node)
	}
	return &Node{Child:child, Edge:edge, suffix:nil, name:name}
}

type active struct {
	n           *Node
	source      []byte
	first, last int 
}

func (a active) explicit() bool { return a.first >= a.last }
func (a active) nextByte() byte {
	return a.n.Child[a.source[a.first]].Edge[a.last-a.first]
}
func (a *active) canonize() {
	for !a.explicit(){
		c := a.source[a.first]
		child := a.n.Child[c]
		if len(child.Edge) > (a.last - a.first) {
			break
		}
		a.n = child
		a.first += len(child.Edge)
	}
}

func (a *active) split() *Node{
	child := a.n.Child[a.source[a.first]]
	child.split(a.last - a.first)
	return child
}

func (n *Node) split(pos int) {
	child := newNode(n.Edge[pos:], n.Child)
	n.Child = map[byte]*Node{child.Edge[0]: child}
	n.Edge = n.Edge[:pos]
}

// New returns the root of a suffix tree for a string.
// It is the caller's responsibility to make sure the last byte
// of the string is unique (otherwise the suffix tree will still have implicit nodes).
func New(s []byte) *Node {
	root := newNode(nil,nil)
	a := active{n: root, source: s}
	for pos := range s {
		a.add(s[pos:], root)
	}
	return root
}

func (a *active) add(suffix []byte, root *Node) {
	var parent, lastParent *Node
	c := suffix[0]
	for {
		parent = a.n
		if a.explicit() {
			if a.n.Child[c] != nil {
				break
			}
		} else {
			if a.nextByte() == c {
				break
			}
			parent = a.split()
		}
		parent.Child[c] = newNode(suffix,nil)
		if lastParent != nil {
			lastParent.suffix = parent
		}
		lastParent = parent

		// move to the next smaller suffix
		if a.n == root {
			a.first++
		} else {
			a.n = a.n.suffix
		}
		a.canonize()
	}

	if lastParent != nil {
		lastParent.suffix = parent
	}
	a.last++ // Now the endpoint is the next active point
	a.canonize()
}
