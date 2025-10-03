package gst_test

import (
	"cmp"
	"fmt"
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/bot/pun/gst"
)

func TestSerialization(t *testing.T) {
	//* === TEST - ALL ===

	for groupNum, testGroup := range tests {
		if !t.Run(fmt.Sprintf("GST-group#%d", groupNum), func(tg *testing.T) {

			//* === TEST - GROUP ===

			var (
				want gst.SuffixTree
				got  gst.SuffixTree
				data []byte
				err1 error
				err2 error
			)

			for testNum, testCase := range testGroup {
				if !tg.Run(fmt.Sprintf("GST-test#%03d", testNum), func(tt *testing.T) {

					//* === TEST - CASE ===

					defer func() {
						if x := recover(); x != nil {
							tt.Fatalf("#%03d fail: PUT panic: %s\nvalue: %s\n", testNum, x, testCase.Value)
						}
					}()

					if testCase.IsPut {
						want.Put([]byte(testCase.Value))
					}

					data, err1 = gst.Serialize(want)
					got, err2 = gst.Deserialize(data)
					if cmp.Or(err1, err2) != nil {
						tt.Fatalf("BIN error: serialization error\n  serialize error: %v\ndeserialize error: %v\n", err1, err2)
					}
					if !reflect.DeepEqual(want, got) {
						tt.Logf("RES error: unexpected result\n\nwant:\n")
						want.Print(tt)
						tt.Logf("\n\n got:\n")
						got.Print(tt)
						tt.FailNow()
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
