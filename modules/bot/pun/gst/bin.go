package gst

func (tree *SuffixTree) AppendBinary(b []byte) ([]byte, error) {
	if tree == nil {
		panic("given suffix tree is nil pointer")
	}
	return nil, nil
}

func (tree *SuffixTree) MarshalBinary() (data []byte, err error) {
	if tree == nil {
		panic("given suffix tree is nil pointer")
	}
	return nil, nil
}

func (tree *SuffixTree) UnmarshalBinary(data []byte) error {
	if tree == nil {
		panic("given suffix tree is nil pointer")
	}
	return nil
}
