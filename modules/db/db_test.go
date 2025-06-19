package db_test

import (
	"log/slog"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/mlog"
)

const PathToRoot string = "../../"

func TestInit(t *testing.T) {
	mlog.InitTesting(t, PathToRoot)
	err := db.InitTesting(t, PathToRoot)
	if err != nil {
		slog.Error("db initialization failed", "error", err)
		t.Fatal(err)
	}
	slog.Info("db initialized correctly, test completed")
}
