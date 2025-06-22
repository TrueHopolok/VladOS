package mlog_test

import (
	"testing"

	"github.com/TrueHopolok/VladOS/modules/mlog"
)

const PathToRoot string = "../../"

func TestInit(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()
	mlog.InitTesting(t, PathToRoot)
}
