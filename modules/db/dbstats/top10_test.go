package dbstats_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
)

func TestTop10(t *testing.T) {
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

	var (
		got  []dbstats.Placement
		want []dbstats.Placement
		err  error
	)
	for _, gameName := range tablesToTest {
		want = nil
		got, err = dbstats.GetTop10(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, _, _, _, _); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Placement{{_}}
		got, err = dbstats.GetTop10(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
	}
}
