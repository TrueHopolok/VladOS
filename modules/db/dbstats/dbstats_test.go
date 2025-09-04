package dbstats_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
)

const pathToRoot = "../../../"

var tablesToTest = []string{"slot", "dice", "bjack", "guess"}

func TestGet(t *testing.T) {
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

	for _, gameName := range tablesToTest {
		if err := dbstats.Update(gameName, 0, 5); err != nil {
			t.Fatal(err)
		}
		stats, err := dbstats.Get(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		want := dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 5, ScoreBest: 5}

		if err := dbstats.Update(gameName, 1, 0); err != nil {
			t.Fatal(err)
		}
		stats, err = dbstats.Get(gameName, 1)
		if err != nil {
			t.Fatal(err)
		}
		want = dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 0, ScoreBest: 0}
		if !reflect.DeepEqual(stats, want) {
			t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
		}

		if err := dbstats.Update(gameName, 1, 1); err != nil {
			t.Fatal(err)
		}
		if err := dbstats.Update(gameName, 1, 64); err != nil {
			t.Fatal(err)
		}
		stats, err = dbstats.Get(gameName, 1)
		if err != nil {
			t.Fatal(err)
		}
		want = dbstats.UserStats{GamesTotal: 3, ScoreCurrent: 65, ScoreBest: 65}
		if !reflect.DeepEqual(stats, want) {
			t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
		}

		if err := dbstats.Update(gameName, 1, 0); err != nil {
			t.Fatal(err)
		}
		if err := dbstats.Update(gameName, 1, 6); err != nil {
			t.Fatal(err)
		}
		stats, err = dbstats.Get(gameName, 1)
		if err != nil {
			t.Fatal(err)
		}
		want = dbstats.UserStats{GamesTotal: 5, ScoreCurrent: 6, ScoreBest: 65}
		if !reflect.DeepEqual(stats, want) {
			t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
		}
	}
}

func TestFull(t *testing.T) {
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
		stats []dbstats.FullStats
		err   error
		want  []dbstats.FullStats
	)

	for _, gameName := range tablesToTest {
		// None players test
		stats, err = dbstats.GetFull(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if stats != nil || len(stats) != 0 {
			t.Fatal("unexpected length of stats")
		}

		// player 1 added
		if err := dbstats.Update(gameName, 1, 1); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.FullStats{
			{UserId: 1, Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 1, ScoreBest: 1}, Placement: 1, PlayersTotal: 1},
		}
		// New player
		stats, err = dbstats.GetFull(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(stats, want) {
			t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
		}
		// Old player
		stats, err = dbstats.GetFull(gameName, 1)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(stats, want) {
			t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
		}

		// Top3 is filled
		if err := dbstats.Update(gameName, 1, 1); err != nil {
			t.Fatal(err)
		}
		if err := dbstats.Update(gameName, 2, 3); err != nil {
			t.Fatal(err)
		}
		if err := dbstats.Update(gameName, 3, 1); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.FullStats{
			{UserId: 2, Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 3, ScoreBest: 3}, Placement: 1, PlayersTotal: 3},
			{UserId: 1, Personal: dbstats.UserStats{GamesTotal: 2, ScoreCurrent: 2, ScoreBest: 2}, Placement: 2, PlayersTotal: 3},
			{UserId: 3, Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 1, ScoreBest: 1}, Placement: 3, PlayersTotal: 3},
		}
		// New player
		stats, err = dbstats.GetFull(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(stats, want) {
			t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
		}
		// Old player
		stats, err = dbstats.GetFull(gameName, 1)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(stats, want) {
			t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
		}

		// Top5 is filled
		if err := dbstats.Update(gameName, 4, 5); err != nil {
			t.Fatal(err)
		}
		if err := dbstats.Update(gameName, 5, 0); err != nil {
			t.Fatal(err)
		}
		want = []dbstats.FullStats{
			{UserId: 4, Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 5, ScoreBest: 5}, Placement: 1, PlayersTotal: 5},
			{UserId: 2, Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 3, ScoreBest: 3}, Placement: 2, PlayersTotal: 5},
			{UserId: 1, Personal: dbstats.UserStats{GamesTotal: 2, ScoreCurrent: 2, ScoreBest: 2}, Placement: 3, PlayersTotal: 5},
		}
		// New player
		stats, err = dbstats.GetFull(gameName, 0)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(stats, want) {
			t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
		}
		// Old player in top3
		stats, err = dbstats.GetFull(gameName, 1)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(stats, want) {
			t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
		}
		// Old player outside top3
		want = append(want, dbstats.FullStats{UserId: 5, Personal: dbstats.UserStats{GamesTotal: 1, ScoreCurrent: 0, ScoreBest: 0}, Placement: 5, PlayersTotal: 5})
		stats, err = dbstats.GetFull(gameName, 5)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(stats, want) {
			t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
		}
	}
}

func TestLeaderboard(t *testing.T) {
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

	got, err = dbstats.Leaderboard(tablesToTest[0])
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatal("gotten result is not nil")
	}

	err = dbstats.Update(tablesToTest[0], 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	want = []dbstats.Placement{{1, 1}}
	got, err = dbstats.Leaderboard(tablesToTest[0])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("gotten result is unexpected\ngot: %v\nwant: %v", got, want)
	}

	err = dbstats.Update(tablesToTest[0], 3, 2)
	if err != nil {
		t.Fatal(err)
	}
	err = dbstats.Update(tablesToTest[0], 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	want = []dbstats.Placement{{2, 2}}
	got, err = dbstats.Leaderboard(tablesToTest[0])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("gotten result is unexpected\ngot: %v\nwant: %v", got, want)
	}

	err = dbstats.Update(tablesToTest[0], 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	err = dbstats.Update(tablesToTest[0], 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	err = dbstats.Update(tablesToTest[0], 3, 1)
	if err != nil {
		t.Fatal(err)
	}
	want = []dbstats.Placement{{3, 1}, {2, 1}, {0, 1}}
	got, err = dbstats.Leaderboard(tablesToTest[0])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("gotten result is unexpected\ngot: %v\nwant: %v", got, want)
	}
}
