package gst

import "testing"

func Equal(t *testing.T, a, b SuffixTree) bool {
	if !testing.Testing() {
		panic("Equal function for GST was called outside of testing build")
	}

	return a.root.equal(b.root)
}

func (a *Node) equal(b *Node) bool {
	if a == nil && b == nil {
		return true
	}
}
