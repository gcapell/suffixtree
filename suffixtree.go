// Ukkonen's algorithm to form suffixtree in linear time.
// Based mostly on http://marknelson.us/1996/08/01/suffix-trees/

package suffixtree

import "log"

type Node struct {
	Child map[byte]*Node
	Edge []byte
	suffix *Node
}
func newNode(edge []byte) *Node { return &Node{make(map[byte]*Node), edge, nil}}

type active struct {
	n *Node
	source []byte
	first, last int // index into source string
}
func (a active) explicit() bool { return a.first >= a.last}
func (a active) nextByte() byte {return a.n.Child[a.source[a.first]].Edge[a.last - a.first]}
func (a *active) canonize() {
	log.Printf("canonize(%v)", *a)
	if a.explicit() {
		return
	}
	for {
		c := a.source[a.first]
		child := a.n.Child[c]
		if child == nil {
			log.Fatalf("expected child of %v: %q (%v), %d", *a.n, string(a.source), c, a.first)
		}
		if len(child.Edge) > (a.last - a.first) {
			break
		}
		a.n = child
		a.first += len(child.Edge)
	}
}

func (n *Node) split(pos int) {
	child := &Node{n.Child, n.Edge[pos:], nil}
	n.Child = map[byte]*Node {child.Edge[0]:child}
	n.Edge = n.Edge[:pos]
}

func New(s []byte)*Node {
	root := newNode(nil)
	a := active{n:root, source:s}	
	for pos := range s {
		a.add(s[pos:], root)
	}
	return root
}

func (a *active) add(suffix []byte, root *Node) {
	log.Println("adding suffix", string(suffix))
    var parent, lastParent *Node
	c := suffix[0]
    for  {
        parent = a.n
        if a.explicit() {
			if a.n.Child[c] != nil {
				break
			}
		} else {
			if a.nextByte() == c {
				break
			}
			a.n.split(a.last - a.first)
        }
		log.Printf("adding child %s(%v) to %v", string(c), c, parent)
		parent.Child[c] = newNode(suffix)
        if  lastParent != nil {
			lastParent.suffix = parent
		}
        lastParent = parent

		// move to the next smaller suffix
        if ( a.n == root) {
			a.first++
		} else {
            a.n = a.n.suffix
		}
        a.canonize()
    }

    if ( lastParent != nil ) {
        lastParent.suffix = parent
	}
    a.last++  // Now the endpoint is the next active point
    a.canonize()
}