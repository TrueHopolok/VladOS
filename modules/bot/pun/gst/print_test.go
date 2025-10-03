package gst

import (
	"testing"
)

func (tree *SuffixTree) Print(t *testing.T) {
	if !testing.Testing() {
		panic("Print function for GST was called outside of testing build")
	}
	t.Logf("GST will be printed in DFS order\n")
	print(t, tree.root, 0)
}

func print(t *testing.T, node *Node, level int) {
	indent := ""
	for range level {
		indent += "\t"
	}
	if node == nil {
		t.Logf("%s%d - nil\n", indent, level)
	}
	t.Logf("%s%d - %p - %v\n", indent, level, node, node.Valid)
	level++
	indent += "\t"
	for char, edge := range node.Edges {
		if edge != nil {
			path := string(char)
			for _, part := range edge.Path {
				path += string(part)
			}
			t.Logf("%s%d - %s\n", indent, level, path)
			print(t, edge.Dest, level)
		}
	}
}
