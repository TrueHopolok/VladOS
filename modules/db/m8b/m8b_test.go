package m8b_test

import (
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/m8b"
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
	if err := db.Migrate(); err != nil {
		t.Fatal(err)
	}

	if _, err := m8b.Get(true); err != nil {
		t.Fatal(err)
	}
	if _, err := m8b.Get(false); err != nil {
		t.Fatal(err)
	}
}
