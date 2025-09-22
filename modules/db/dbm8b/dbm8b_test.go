package dbm8b_test

import (
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbm8b"
)

const pathToRoot = "../../../"

func TestM8B(t *testing.T) {
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

	if _, err := dbm8b.Get(true); err != nil {
		t.Fatal(err)
	}
	if _, err := dbm8b.Get(false); err != nil {
		t.Fatal(err)
	}

	if err := dbm8b.Add(dbm8b.M8B{Answer: true, Text: "ahah"}); err != nil {
		t.Fatal(err)
	}
}
