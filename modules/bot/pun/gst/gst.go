// GST package is realization of generalized suffixtree for the VladOS project.
//
// Realization is lazy, since the depth of any given tree is 3 characters based on project's limitation.
// In addition, only lower case english letters are considered.
//
// Based on that, max size of any given tree is 26^3 or 17576.
// This allows for even a full tree to be easily stored in database.
package gst

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

type Node struct {
	// Valid indicates whether or not currently given suffix is valid.
	//
	// If is invalid => suffix at this node does not exists => children may contain existing suffix.
	Valid bool

	// Index stores id of an entry in database for faster search in it.
	//
	// This value can be used if [Valid] is true.
	Index int

	// Contain all information regarding the
	Edges [26]*Edge
}

type Edge struct {
	// Path is symbols that are required to reach a node.
	Path string

	// Dest is destination of the Edge.
	//
	// Dest should not be nil.
	Dest *Node
}

type SuffixTree struct {
	// Size is amount of strings that were put in to the tree.
	//
	// Based on project limitation, max size if 26^3 or 17576.
	Size int

	// Root is starting node of any tree.
	//
	// If tree is empty, root will be nil.
	Root *Node
}
