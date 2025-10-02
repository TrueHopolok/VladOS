package gst

// Get for a given word return the longest suffix that is stored in the tree.
//
// Panics if tree is nil.
func (tree *SuffixTree) Get(word []byte) (maxSuffix []byte) {
	if tree == nil {
		panic("given suffix tree is nil pointer")
	}
	wordIndex := len(word) - 1
	if tree.root == nil || wordIndex < 0 {
		return
	}

	maxSuffix = make([]byte, 0)
	suffix := make([]byte, 0)
	current := tree.root

	for {
		if wordIndex < 0 {
			return
		}

		edge := current.Edges[word[wordIndex]]
		if edge == nil {
			return
		}

		suffix = append(suffix, word[wordIndex])
		wordIndex--

		for edgeIndex := 0; edgeIndex < len(edge.Path); edgeIndex++ {

			if wordIndex < 0 {
				return
			}

			if edge.Path[edgeIndex] != word[wordIndex] {
				return
			}

			suffix = append(suffix, word[wordIndex])
			wordIndex--
		}

		current = edge.Dest
		if current == nil {
			panic("impossible scenario: destination of the edge is nil")
		}
		if current.Valid {
			maxSuffix = append(maxSuffix, suffix[len(maxSuffix):]...)
		}
	}
}
