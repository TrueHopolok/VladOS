package db_test

import (
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/mlog"
)

const pathToRoot string = "../../"

func TestInit(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()
	err := db.InitTesting(t, pathToRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Conn.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestMigrate(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()
	mlog.InitTesting(t, pathToRoot) // used since mlog logs its actions
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
	if err := db.Migrate(); err != nil {
		t.Fatal(err)
	}
}
