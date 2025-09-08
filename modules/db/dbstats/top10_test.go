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

		if err = dbstats.Update(gameName, 0, "0", "", 0); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Placement{
			{UserId: 0, FirstName: "0", Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 0, ScoreBest: 0}, Placement: 1, PlayersTotal: 1},
		}
		got, err = dbstats.GetTop10(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 0, "0", "username", 12); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Placement{
			{UserId: 0, FirstName: "0", Username: "username", Personal: dbstats.UserStats{GamesTotal: 2, ScoreCurrent: 12, ScoreBest: 12}, Placement: 1, PlayersTotal: 1},
		}
		got, err = dbstats.GetTop10(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 1, "1", "", 5); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Placement{
			{UserId: 0, FirstName: "0", Username: "username", Personal: dbstats.UserStats{GamesTotal: 2, ScoreCurrent: 12, ScoreBest: 12}, Placement: 1, PlayersTotal: 2},
			{UserId: 1, FirstName: "1", Username: "", Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 5, ScoreBest: 5}, Placement: 2, PlayersTotal: 2},
		}
		got, err = dbstats.GetTop10(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 0, "0", "0", 0); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Placement{
			{UserId: 0, FirstName: "0", Username: "0", Personal: dbstats.UserStats{GamesTotal: 3, ScoreCurrent: 0, ScoreBest: 12}, Placement: 1, PlayersTotal: 2},
			{UserId: 1, FirstName: "1", Username: "", Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 5, ScoreBest: 5}, Placement: 2, PlayersTotal: 2},
		}
		got, err = dbstats.GetTop10(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 1, "12", "", 12); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Placement{
			{UserId: 1, FirstName: "12", Username: "", Personal: dbstats.UserStats{GamesTotal: 2, ScoreCurrent: 17, ScoreBest: 17}, Placement: 1, PlayersTotal: 2},
			{UserId: 0, FirstName: "0", Username: "0", Personal: dbstats.UserStats{GamesTotal: 3, ScoreCurrent: 0, ScoreBest: 12}, Placement: 2, PlayersTotal: 2},
		}
		got, err = dbstats.GetTop10(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		want = []dbstats.Placement{
			{UserId: 1, FirstName: "12", Username: "", Personal: dbstats.UserStats{GamesTotal: 2, ScoreCurrent: 17, ScoreBest: 17}, Placement: 1, PlayersTotal: 12},
			{UserId: 0, FirstName: "0", Username: "0", Personal: dbstats.UserStats{GamesTotal: 3, ScoreCurrent: 0, ScoreBest: 12}, Placement: 2, PlayersTotal: 12},
		}
		for i := 11; i >= 2; i-- {
			if err = dbstats.Update(gameName, int64(i), "a", "", i); err != nil {
				t.Fatal(err)
			}
			if len(want) < 10 {
				want = append(want, dbstats.Placement{UserId: int64(i), FirstName: "a", Username: "", Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: i, ScoreBest: i}, Placement: 14 - i, PlayersTotal: 12})
			}
		}
		got, err = dbstats.GetTop10(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
	}
}
