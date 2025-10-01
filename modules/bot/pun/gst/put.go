package gst

import "errors"

func (tree *SuffixTree) Put(word string) error {
	if tree == nil {
		panic("given suffix tree is nil pointer")
	}
	if len(word) == 0 {
		return errors.New("word is empty")
	}
	wordIndex := len(word) - 1

	// Project limitations
	if wordIndex > 2 {
		return errors.New("word is too long")
	}
	for _, c := range word {
		c -= 'a'
		if c > 25 || c < 0 {
			return errors.New("word contain invalid symbols")
		}
	}

	if tree.Size == 0 {
		tree.Size++
		tree.Root = &Node{}
	}

	current := tree.Root
	for {
		// TODO
		if current.Edges[word[wordIndex]] == nil {
			edge := &Edge{}
			current.Edges[word[wordIndex]] = edge
			for wordIndex--; wordIndex >= 0; wordIndex-- {
				c := word[wordIndex]
				current.Edges[o].Path += c
			}
		}
	}
}
