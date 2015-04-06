package suffixtree

import (
	"io"
	"io/ioutil"
	"os"
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func digraph(root *Node, a active, o io.Writer) {
	fmt.Fprintln(o, "digraph G {")
	root.graphChildren(o)
	root.graphSuffixes(o)
	fmt.Fprintf(o, "active -> %s [label=\"%s\"]", nodeName(a.n), string(a.source[a.first:a.last]))
	fmt.Fprintln(o, "}")
}

func png(root *Node, a active, filename string)error {
	dot, err := exec.LookPath("dot")
	if err != nil {
		return err
	}
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	o := bufio.NewWriter(f)
	digraph(root, a, o)
	o.Flush()
	f.Close()
	
	commandLine := []string{dot, "-Tpng", f.Name(), "-o", filename}
	cmd := exec.Command(commandLine[0], commandLine[1:]...)
	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s executing %q\n", err, commandLine)
		return err
	}
	os.Remove(f.Name())
	return nil
}

func (n *Node)graphChildren(o io.Writer) {
	name := nodeName(n)
	for _, child := range n.Child {
		fmt.Fprintf(o, "%s -> %s [ label=\"%s\"];\n", name, nodeName(child), string(child.Edge))
	}
	for _, child := range n.Child {
		child.graphChildren(o)
	}
}

func (n *Node)graphSuffixes(o io.Writer) {
	if n.suffix != nil {
		fmt.Fprintf(o, "%s -> %s [ style=dotted];\n", nodeName(n), nodeName(n.suffix))
	}
		for _, child := range n.Child {
		child.graphSuffixes(o)
	}
}

var (
	nodeNames = make(map[*Node]string)
	lastNodeName int
)

func nodeName(n *Node)string {
	if name, ok := nodeNames[n]; ok {
		return name
	}
	name := strconv.Itoa(lastNodeName)
	lastNodeName++
	nodeNames[n] = name
	return name
}

func diNode(edge string, child... *Node)*Node{
	n := &Node{make(map[byte]*Node), []byte(edge), nil}
	for _, c := range child {
		n.Child[c.Edge[0]] = c
	}
	return n
}

func BuildAndGraph(s []byte) {
	root := newNode(nil)
	a := active{n:root, source:s}
	if err := png(root, a, "orig.png"); err != nil {
		log.Fatal(err)
	}	
	for pos := range s {
		a.add(s[pos:], root)
		if err := png(root, a, fmt.Sprintf("%d.png", pos)); err != nil {
			log.Fatal(err)
		}
	}	
}

