package dbstats_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
)

func TestTopSelf(t *testing.T) {
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
		got, err = dbstats.GetTopSelf(gameName, 0)
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
		got, err = dbstats.GetTopSelf(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
		got, err = dbstats.GetTopSelf(gameName, 1)
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
		got, err = dbstats.GetTopSelf(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
		got, err = dbstats.GetTopSelf(gameName, 1)
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
		got, err = dbstats.GetTopSelf(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
		got, err = dbstats.GetTopSelf(gameName, 2)
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
		got, err = dbstats.GetTopSelf(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
		got, err = dbstats.GetTopSelf(gameName, 2)
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
		got, err = dbstats.GetTopSelf(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
		got, err = dbstats.GetTopSelf(gameName, 2)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}

		if err = dbstats.Update(gameName, 2, "2", "", 0); err != nil {
			t.Fatal(err)
		}
		if err = dbstats.Update(gameName, 3, "3", "", 5); err != nil {
			t.Fatal(err)
		}
		if err = dbstats.Update(gameName, 4, "4", "", 15); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.Placement{
			{UserId: 1, FirstName: "12", Username: "", Personal: dbstats.UserStats{GamesTotal: 2, ScoreCurrent: 17, ScoreBest: 17}, Placement: 1, PlayersTotal: 5},
			{UserId: 4, FirstName: "4", Username: "", Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 15, ScoreBest: 15}, Placement: 2, PlayersTotal: 5},
			{UserId: 0, FirstName: "0", Username: "0", Personal: dbstats.UserStats{GamesTotal: 3, ScoreCurrent: 0, ScoreBest: 12}, Placement: 3, PlayersTotal: 5},
		}
		got, err = dbstats.GetTopSelf(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
		got, err = dbstats.GetTopSelf(gameName, 6)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
		want = append(want, dbstats.Placement{UserId: 2, FirstName: "2", Username: "", Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 0, ScoreBest: 0}, Placement: 5, PlayersTotal: 5})
		got, err = dbstats.GetTopSelf(gameName, 4)
		if err != nil {
			t.Fatal(err)
		}
	}
}
