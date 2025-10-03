package gst

import (
	"bytes"
	"encoding/gob"
)

func Serialize(tree SuffixTree) ([]byte, error) {
	if tree.root == nil {
		return nil, nil
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(tree.root)
	return buf.Bytes(), err
}

func Deserialize(data []byte) (SuffixTree, error) {
	var tree SuffixTree
	if len(data) == 0 {
		return tree, nil
	}

	var buf bytes.Buffer
	if _, err := buf.Write(data); err != nil {
		return tree, err
	}
	dec := gob.NewDecoder(&buf)
	err := dec.Decode(&tree.root)
	return tree, err
}
