package gst_test

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

var tests = [][]Case{
	{
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
	},
	{
		{IsPut: true, Value: "a"},
		{IsPut: true, Value: "ja"},
		{IsPut: true, Value: "il"},
		{IsPut: true, Value: "ti"},
		{IsPut: true, Value: "eayo"},

		{IsPut: false, Value: "privet", Expected: ""},
		{IsPut: false, Value: "ti sprosil menja", Expected: "ja"},
		{IsPut: false, Value: "gde ti", Expected: "ti"},
		{IsPut: false, Value: "yo yo yo dorogije podpisiki", Expected: ""},
		{IsPut: false, Value: "poka", Expected: "a"},
	},
}
