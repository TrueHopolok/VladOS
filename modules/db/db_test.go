package db_test

import (
	"log/slog"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
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
	slog.Info("db.TestInit", "STATUS", "START")
	err := db.InitTesting(t, PathToRoot)
	if err != nil {
		slog.Error("db.TestInit", "STATUS", "FAILED", "error", err)
		t.Fatal(err)
	}
	slog.Info("db.TestInit", "STATUS", "SUCCESS")
}

func TestMigrate(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()
	mlog.InitTesting(t, PathToRoot)
	slog.Info("db.TestMigrate", "STATUS", "START")
	if err := db.InitTesting(t, PathToRoot); err != nil {
		slog.Info("db.TestMigrate", "STATUS", "FAILED", "error", err)
		t.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		slog.Info("db.TestMigrate", "STATUS", "FAILED", "error", err)
		t.Fatal(err)
	}
	if err := db.Migrate(); err != nil {
		slog.Info("db.TestMigrate", "STATUS", "FAILED", "error", err)
		t.Fatal(err)
	}
	slog.Info("db.TestMigrate", "STATUS", "SUCCESS")
}
