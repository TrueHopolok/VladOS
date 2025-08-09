package dbtip_test

import (
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbtip"
)

const pathToRoot = "../../../"

func TestTip(t *testing.T) {
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

	_, _, id, err := dbtip.Rand()
	if err != nil {
		t.Fatal(err)
	}

	_, _, found, err := dbtip.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("tip is not found in Get() even though it was found in Rand()")
	}
}
