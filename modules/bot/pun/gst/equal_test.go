package gst

import "testing"

func Equal(t *testing.T, a, b SuffixTree) bool {
	if !testing.Testing() {
		panic("Equal function for GST was called outside of testing build")
	}

	return a.root.equal(b.root)
}

func (a *Node) equal(b *Node) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.Valid != b.Valid {
		return false
	}
	if len(a.Edges) != len(b.Edges) {
		return false
	}
	for key := range a.Edges {
		if !a.Edges[key].equal(b.Edges[key]) {
			return false
		}
	}
	return true
}

func (a *Edge) equal(b *Edge) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a.Path) != len(b.Path) {
		return false
	}
	for i := range a.Path {
		if a.Path[i] != b.Path[i] {
			return false
		}
	}
	return a.Dest.equal(b.Dest)
}
