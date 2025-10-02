package gst

import (
	"errors"
)

func (tree *SuffixTree) Put(word []byte) error {
	if tree == nil {
		panic("given suffix tree is nil pointer")
	}
	if len(word) == 0 {
		return errors.New("word is empty")
	}
	wordIndex := len(word) - 1

	// Project limitations
	for i := range word {
		word[i] -= 'a'
		if word[i] > 25 {
			return errors.New("word contain invalid symbols")
		}
	}

	if tree.root == nil {
		tree.root = &Node{}
	}

	current := tree.root
	for {
		if wordIndex < 0 {
			/*
				node.valid = true
			*/
			current.Valid = true
			return nil
		}

		edge := current.Edges[word[wordIndex]]
		if edge == nil {
			/*
				old_node -> new_edge -> new_node
			*/
			edge := &Edge{}
			current.Edges[word[wordIndex]] = edge
			for wordIndex--; wordIndex >= 0; wordIndex-- {
				edge.Path = append(edge.Path, word[wordIndex])
			}
			edge.Dest = &Node{Valid: true}
			return nil
		}

		wordIndex--
		for edgeIndex := 0; edgeIndex < len(edge.Path); edgeIndex++ {

			if wordIndex < 0 {
				/*
					original_edge -> new_node
					new_node.valid = true
					new_node -> new_edge + part of the old path
				*/
				oldPath := edge.Path[edgeIndex:]
				oldNode := edge.Dest
				edge.Path = edge.Path[:edgeIndex]
				edge.Dest = &Node{Valid: true}
				edge.Dest.Edges[oldPath[0]] = &Edge{Path: oldPath[1:], Dest: oldNode}
				return nil
			}

			if word[wordIndex] != edge.Path[edgeIndex] {
				/*
					original_edge -> new_node
					new_node -> old_path -> old_node
					new_node -> new_path -> new_node
				*/
				oldPath := edge.Path[edgeIndex:]
				oldNode := edge.Dest
				edge.Path = edge.Path[:edgeIndex]
				edge.Dest = &Node{}
				edge.Dest.Edges[oldPath[0]] = &Edge{Path: oldPath[1:], Dest: oldNode}
				newEdge := &Edge{Dest: &Node{Valid: true}}
				edge.Dest.Edges[word[wordIndex]] = newEdge
				for wordIndex--; wordIndex >= 0; wordIndex-- {
					newEdge.Path = append(newEdge.Path, word[wordIndex])
				}
				return nil
			}

			wordIndex--
		}

		current = edge.Dest
	}
}
