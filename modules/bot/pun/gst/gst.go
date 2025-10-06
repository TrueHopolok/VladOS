// GST package is realization of generalized suffixtree for the VladOS project.
// This package supports trees of any depth and any ASCII symbols / strings.
//
// Realization is based on my knowledge and intuition (might not be the fastest),
// but projet limitations are small, thus not requiring a fast implementation.
package gst

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

type Node struct {
	// Valid indicates whether or not currently given suffix is valid.
	//
	// If is invalid => suffix at this node does not exists => children may contain existing suffix.
	Valid bool

	// Contain all information regarding the edges out goind out of a node.
	Edges map[byte]*Edge
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
