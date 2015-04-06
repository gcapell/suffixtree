package suffixtree

import (
	"io"
	"io/ioutil"
	"os"
	"bufio"
	"fmt"
	"os/exec"
	"log"
)

func digraph(root *Node, a active, o io.Writer) {
	fmt.Fprintln(o, "digraph G {")
	root.graphChildren(o)
	root.graphSuffixes(o)
	fmt.Fprintf(o, "active -> %s [label=\"%s\"]", a.n.name, string(a.source[a.first:a.last]))
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
	for _, child := range n.Child {
		fmt.Fprintf(o, "%s -> %s [ label=\"%s\"];\n", n.name, child.name, string(child.Edge))
	}
	for _, child := range n.Child {
		child.graphChildren(o)
	}
}

func (n *Node)graphSuffixes(o io.Writer) {
	if n.suffix != nil {
		fmt.Fprintf(o, "%s -> %s [ style=dotted];\n", n.name, n.suffix.name)
	}
		for _, child := range n.Child {
		child.graphSuffixes(o)
	}
}

func BuildAndGraph(s []byte) {
	root := newNode(nil,nil)
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

