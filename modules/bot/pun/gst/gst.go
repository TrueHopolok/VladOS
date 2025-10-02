// GST package is realization of generalized suffixtree for the VladOS project.
// This package supports trees of any depth, with single limitation of using only lower case english letters.
//
// Realization is based on my knowledge and intuition (might not be the fastest),
// but projet limitations are small, thus not requiring a fast implementation.
//
// Max suffix is the size of 3, thus max size of any given tree is 26^3 or 17576.
// This allows for even a full tree to be easily stored in database.
// For that [SuffixTree] implements binary marshal / unmarshal / appender interfaces.
package gst

import "fmt"

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

type Node struct {
	// Valid indicates whether or not currently given suffix is valid.
	//
	// If is invalid => suffix at this node does not exists => children may contain existing suffix.
	Valid bool

	// Contain all information regarding the
	Edges [26]*Edge
}

type Edge struct {
	// Path is symbols that are required to reach a node.
	Path []byte

	// Dest is destination of the Edge.
	//
	// Dest should not be nil.
	Dest *Node
}

type SuffixTree struct {
	// root is starting node of any tree.
	//
	// If tree is empty, root will be nil.
	root *Node
}

func (tree *SuffixTree) Print() {
	if tree == nil {
		panic("given suffix tree is nil pointer")
	}
	print(tree.root, 1)
}

func print(node *Node, level int) {
	if node == nil {
		fmt.Printf("%d - nil\n", level)
	}
	fmt.Printf("%d - %v\n", level, node.Valid)
	level++
	for char, edge := range node.Edges {
		if edge != nil {
			fmt.Printf("%s", string(byte(char)+'a'))
			for _, path := range edge.Path {
				fmt.Printf("%s", string(byte(path)+'a'))
			}
			fmt.Print(" | ")
			print(edge.Dest, level)
		}
	}
}
