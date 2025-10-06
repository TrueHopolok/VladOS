package dbpun_test

import (
	"errors"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbpun"
)

const pathToRoot = "../../../"

type Case struct {
	isw bool   // is writing operation
	err error  // expected error
	key string // suffix for write or word for answer
	pun string // inserted for write or expected for answer
}

var tests = []Case{
	{
		isw: true,
		err: &dbpun.ArgumentLengthError{"given suffix is empty"},
		key: "",
		pun: "",
	},
	{
		isw: false,
		err: &dbpun.ArgumentLengthError{"given pun is empty"},
		key: "",
		pun: "",
	},
	{
		isw: false,
		err: &dbpun.ArgumentLengthError{"given word is empty"},
		key: "",
		pun: "",
	},
	{
		isw: false,
		err: nil,
		key: "hello",
		pun: "",
	},
	{
		isw: true,
		err: nil,
		key: "o",
		pun: "o",
	},
	{
		isw: false,
		err: nil,
		key: "a",
		pun: "",
	},
	{
		isw: false,
		err: nil,
		key: "o",
		pun: "o",
	},
	{
		isw: false,
		err: nil,
		key: "lo",
		pun: "o",
	},
	{
		isw: false,
		err: nil,
		key: "hello",
		pun: "o",
	},
	{
		isw: false,
		err: nil,
		key: "_hello",
		pun: "o",
	},
	{
		isw: true,
		err: nil,
		key: "hello",
		pun: "hello",
	},
	{
		isw: false,
		err: nil,
		key: "a",
		pun: "",
	},
	{
		isw: false,
		err: nil,
		key: "o",
		pun: "o",
	},
	{
		isw: false,
		err: nil,
		key: "lo",
		pun: "o",
	},
	{
		isw: false,
		err: nil,
		key: "hello",
		pun: "hello",
	},
	{
		isw: false,
		err: nil,
		key: "_hello",
		pun: "hello",
	},
}

func TestPun(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()
	if err := db.InitTesting(t, pathToRoot); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Conn.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	if err := db.Migrate(); err != nil {
		t.Fatal(err)
	}

	var (
		err error
		pun string
	)
	for testNum, testCase := range tests {
		if testCase.isw {
			err = dbpun.Write(testCase.key, testCase.pun)
			if errors.Is(err, &dbpun.ArgumentLengthError{}) && errors.Is(testCase.err, &dbpun.ArgumentLengthError{}) {
				t.Logf("test #%03d - FAIL\n", testNum)
				t.Logf("      type: WRITE\n")
				t.Logf("      diff: error\n")
				t.Logf("       got: %v\n", err)
				t.Logf("      want: %v\n", testCase.err)
				t.Fail()
			}
		} else {
			pun, err = dbpun.Answer(testCase.key)
			if errors.Is(err, &dbpun.ArgumentLengthError{}) && errors.Is(testCase.err, &dbpun.ArgumentLengthError{}) {
				t.Logf("test #%03d - FAIL\n", testNum)
				t.Logf("      type: ANSWER\n")
				t.Logf("      diff: error\n")
				t.Logf("       got: %v\n", err)
				t.Logf("      want: %v\n", testCase.err)
				t.Fail()
			}
			if pun != testCase.pun {
				t.Logf("test #%03d - FAIL\ntype: WRITE\n", testNum)
				t.Logf("      type: ANSWER\n")
				t.Logf("      diff: pun\n")
				t.Logf("       got: %v\n", pun)
				t.Logf("      want: %v\n", testCase.pun)
				t.Fail()
			}
		}
		if t.Failed() {
			t.Fatalf("whole test: %+v", testCase)
		} else {
			t.Logf("test #%03d - pass\n", testNum)
		}
	}
}
