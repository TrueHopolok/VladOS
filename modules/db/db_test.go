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

func TestMigrate(t *testing.T) {
	mlog.InitTesting(t, PathToRoot)
	if err := db.InitTesting(t, PathToRoot); err != nil {
		slog.Error("db migration failed", "error", err)
		t.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		slog.Error("db migration failed", "error", err)
		t.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		slog.Error("db migration failed", "error", err)
		t.Fatal(err)
	}
	slog.Info("db migrted correctly, test completed")
}
