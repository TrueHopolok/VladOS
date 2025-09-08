package dbstats_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
)

func TestPrecent(t *testing.T) {
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
		got  []dbstats.Precent
		want []dbstats.Precent
		err  error
	)
	for _, gameName := range tablesToTest {
		want = nil
		got, err = dbstats.GetPrecent(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 0, "0", "", 0); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Precent{{ScoreBest: 0, PlayersAmount: 1}}
		got, err = dbstats.GetPrecent(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 0, "0", "", 1); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Precent{{ScoreBest: 1, PlayersAmount: 1}}
		got, err = dbstats.GetPrecent(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 1, "1", "", 5); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Precent{{ScoreBest: 5, PlayersAmount: 1}, {ScoreBest: 1, PlayersAmount: 1}}
		got, err = dbstats.GetPrecent(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 0, "0", "", 4); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Precent{{ScoreBest: 5, PlayersAmount: 2}}
		got, err = dbstats.GetPrecent(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 0, "0", "", 0); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Precent{{ScoreBest: 5, PlayersAmount: 2}}
		got, err = dbstats.GetPrecent(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 0, "0", "", 10); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Precent{{ScoreBest: 10, PlayersAmount: 1}, {ScoreBest: 5, PlayersAmount: 1}}
		got, err = dbstats.GetPrecent(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
	}
}
