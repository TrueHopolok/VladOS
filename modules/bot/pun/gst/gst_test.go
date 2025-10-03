package gst_test

import (
	"fmt"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/bot/pun/gst"
)

func TestFunctional(t *testing.T) {
	//* === TEST - ALL ===

	for groupNum, testGroup := range tests {
		if !t.Run(fmt.Sprintf("GST-group#%d", groupNum), func(tg *testing.T) {

			//* === TEST - GROUP ===

			var (
				tree gst.SuffixTree
				got  []byte
			)
			for testNum, testCase := range testGroup {
				if !tg.Run(fmt.Sprintf("GST-test#%03d", testNum), func(tt *testing.T) {

					//* === TEST - CASE ===

					defer func() {
						if x := recover(); x != nil {
							testType := "GET"
							if testCase.IsPut {
								testType = "PUT"
							}
							tt.Fatalf("#%03d fail: %s panic: %s\nvalue: %s\n", testNum, testType, x, testCase.Value)
						}
					}()

					if testCase.IsPut {
						tree.Put([]byte(testCase.Value))
					} else {
						got = tree.Get([]byte(testCase.Value))
						if string(got) != testCase.Expected {
							tt.Errorf("#%03d fail: GET error: unexpected value\n  want: %s\n   got: %s\nsearch: %s\n", testNum, testCase.Expected, got, testCase.Value)
							tree.Print(tt)
							tt.FailNow()
						}
					}

					//! === TEST - CASE ===

				}) {
					tg.FailNow()
				}

				//! === TEST - GROUP ===
			}
		}) {
			t.Fail()
		}
	}

	//! === TEST - ALL ===
}
