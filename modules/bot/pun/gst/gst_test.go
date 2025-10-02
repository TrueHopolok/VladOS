package gst_test

import (
	"testing"

	"github.com/TrueHopolok/VladOS/modules/bot/pun/gst"
)

type Case struct {
	// IsPut decide what type of an test case is that.
	//
	// True means it is Put operation test, while False means it is Get operation test.
	IsPut bool

	// Value is what used during the test as an argument for a Put or Get command.
	// For ease of use, value is stored as string, but then converted into slice of bytes.
	Value string

	// Expected value is what is expected from test case.
	//
	// Field is used only for Get type of the test.
	Expected string
}

func TestGst(t *testing.T) {
	tests := []Case{
		{IsPut: false, Value: "", Expected: ""},
		{IsPut: false, Value: "hello", Expected: ""},

		{IsPut: true, Value: "a"},
		{IsPut: false, Value: "", Expected: ""},
		{IsPut: false, Value: "hello", Expected: ""},
		{IsPut: false, Value: "a", Expected: "a"},
		{IsPut: false, Value: "aa", Expected: "a"},

		{IsPut: true, Value: "dcba"},
		{IsPut: false, Value: "", Expected: ""},
		{IsPut: false, Value: "hello", Expected: ""},
		{IsPut: false, Value: "a", Expected: "a"},
		{IsPut: false, Value: "aa", Expected: "a"},
		{IsPut: false, Value: "ab", Expected: ""},
		{IsPut: false, Value: "ba", Expected: "a"},
		{IsPut: false, Value: "abc", Expected: ""},
		{IsPut: false, Value: "cba", Expected: "a"},
		{IsPut: false, Value: "abcd", Expected: ""},
		{IsPut: false, Value: "dcba", Expected: "dcba"},

		{IsPut: true, Value: "ba"},
		{IsPut: false, Value: "", Expected: ""},
		{IsPut: false, Value: "hello", Expected: ""},
		{IsPut: false, Value: "a", Expected: "a"},
		{IsPut: false, Value: "aa", Expected: "a"},
		{IsPut: false, Value: "ab", Expected: ""},
		{IsPut: false, Value: "aba", Expected: "ba"},
		{IsPut: false, Value: "abc", Expected: ""},
		{IsPut: false, Value: "abca", Expected: "a"},
		{IsPut: false, Value: "dcba", Expected: "dcba"},
		{IsPut: false, Value: "abcda", Expected: "a"},

		{IsPut: true, Value: "ecba"},
		{IsPut: false, Value: "", Expected: ""},
		{IsPut: false, Value: "hello", Expected: ""},
		{IsPut: false, Value: "a", Expected: "a"},
		{IsPut: false, Value: "aa", Expected: "a"},
		{IsPut: false, Value: "ab", Expected: ""},
		{IsPut: false, Value: "aba", Expected: "ba"},
		{IsPut: false, Value: "abc", Expected: ""},
		{IsPut: false, Value: "abca", Expected: "a"},
		{IsPut: false, Value: "abcd", Expected: ""},
		{IsPut: false, Value: "abcda", Expected: "a"},
		{IsPut: false, Value: "dcba", Expected: "dcba"},
		{IsPut: false, Value: "ecba", Expected: "ecba"},
		{IsPut: false, Value: "recba", Expected: "ecba"},
	}
	var (
		tree gst.SuffixTree
		got  []byte
		err  error
	)
	defer func() {
		if x := recover(); x != nil {
			t.Fatalf("testing paniced: %v", x)
		}
	}()
	for testNum, testCase := range tests {
		t.Logf("#%03d started\n", testNum)
		if testCase.IsPut {
			err = tree.Put([]byte(testCase.Value))
			if err != nil {
				tree.Print()
				t.Fatalf("#%03d failed: PUT error: %s\nvalue: %s", testNum, err, testCase.Value)
			}
		} else {
			got = tree.Get([]byte(testCase.Value))
			if string(got) != testCase.Expected {
				tree.Print()
				t.Fatalf("#%03d failed: GET error: unexpected value\n  want: %s\n   got: %s\nsearch: %s", testNum, testCase.Expected, got, testCase.Value)
			}
		}
		t.Logf("#%03d ok\n", testNum)
	}
}
