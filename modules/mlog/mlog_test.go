package mlog_test

import (
	"log/slog"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/mlog"
)

const PathToRoot string = "../../"

func TestInit(t *testing.T) {
	mlog.InitTesting(t, PathToRoot)
	slog.Info("logger initialized correctly, test completed")
}
