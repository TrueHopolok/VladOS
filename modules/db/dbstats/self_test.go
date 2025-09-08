package dbstats_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
)

func TestSelf(t *testing.T) {
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
		got  dbstats.UserStats
		want dbstats.UserStats
		err  error
	)
	for _, gameName := range tablesToTest {
		want = dbstats.UserStats{}
		got, err = dbstats.GetSelf(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 0, "0", "", 0); err != nil {
			t.Fatal(err)
		}
		want.GamesTotal = 1
		got, err = dbstats.GetSelf(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 1, "1", "", 5); err != nil {
			t.Fatal(err)
		}
		got, err = dbstats.GetSelf(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
		want = dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 5, ScoreBest: 5}
		got, err = dbstats.GetSelf(gameName, 1)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 1, "1", "", 0); err != nil {
			t.Fatal(err)
		}
		want = dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 0, ScoreBest: 0}
		got, err = dbstats.GetSelf(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
		want = dbstats.UserStats{GamesTotal: 2, ScoreCurrent: 0, ScoreBest: 5}
		got, err = dbstats.GetSelf(gameName, 1)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
	}
}
