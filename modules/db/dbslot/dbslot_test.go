package dbslot_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbslot"
)

const pathToRoot = "../../../"

func equalStats(x, y dbslot.UserStats) bool {
	return x.SpinsTotal == y.SpinsTotal && x.ScoreCurrent == y.ScoreCurrent && x.ScoreBest == y.ScoreBest
}

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

	if err := dbslot.Update(0, 5); err != nil {
		t.Fatal(err)
	}
	stats, err := dbslot.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	want := dbslot.UserStats{SpinsTotal: 1, ScoreCurrent: 5, ScoreBest: 5}

	if err := dbslot.Update(1, 0); err != nil {
		t.Fatal(err)
	}
	stats, err = dbslot.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	want = dbslot.UserStats{SpinsTotal: 1, ScoreCurrent: 0, ScoreBest: 0}
	if !equalStats(stats, want) {
		t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
	}

	if err := dbslot.Update(1, 1); err != nil {
		t.Fatal(err)
	}
	if err := dbslot.Update(1, 64); err != nil {
		t.Fatal(err)
	}
	stats, err = dbslot.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	want = dbslot.UserStats{SpinsTotal: 3, ScoreCurrent: 65, ScoreBest: 65}
	if !equalStats(stats, want) {
		t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
	}

	if err := dbslot.Update(1, 0); err != nil {
		t.Fatal(err)
	}
	if err := dbslot.Update(1, 6); err != nil {
		t.Fatal(err)
	}
	stats, err = dbslot.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	want = dbslot.UserStats{SpinsTotal: 5, ScoreCurrent: 6, ScoreBest: 65}
	if !equalStats(stats, want) {
		t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
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
		stats []dbslot.FullStats
		err   error
		want  []dbslot.FullStats
	)

	// None players test
	stats, err = dbslot.GetFull(0)
	if err != nil {
		t.Fatal(err)
	}
	if stats != nil || len(stats) != 0 {
		t.Fatal("unexpected length of stats")
	}

	// player 1 added
	if err := dbslot.Update(1, 1); err != nil {
		t.Fatal(err)
	}
	want = []dbslot.FullStats{
		{UserId: 1, Personal: dbslot.UserStats{SpinsTotal: 1, ScoreCurrent: 1, ScoreBest: 1}, Placement: 1, PlayersTotal: 1},
	}
	// New player
	stats, err = dbslot.GetFull(0)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(stats, want) {
		t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
	}
	// Old player
	stats, err = dbslot.GetFull(1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(stats, want) {
		t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
	}

	// Top3 is filled
	if err := dbslot.Update(1, 1); err != nil {
		t.Fatal(err)
	}
	if err := dbslot.Update(2, 3); err != nil {
		t.Fatal(err)
	}
	if err := dbslot.Update(3, 1); err != nil {
		t.Fatal(err)
	}
	want = []dbslot.FullStats{
		{UserId: 2, Personal: dbslot.UserStats{SpinsTotal: 1, ScoreCurrent: 3, ScoreBest: 3}, Placement: 1, PlayersTotal: 3},
		{UserId: 1, Personal: dbslot.UserStats{SpinsTotal: 2, ScoreCurrent: 2, ScoreBest: 2}, Placement: 2, PlayersTotal: 3},
		{UserId: 3, Personal: dbslot.UserStats{SpinsTotal: 1, ScoreCurrent: 1, ScoreBest: 1}, Placement: 3, PlayersTotal: 3},
	}
	// New player
	stats, err = dbslot.GetFull(0)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(stats, want) {
		t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
	}
	// Old player
	stats, err = dbslot.GetFull(1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(stats, want) {
		t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
	}

	// Top5 is filled
	if err := dbslot.Update(4, 5); err != nil {
		t.Fatal(err)
	}
	if err := dbslot.Update(5, 0); err != nil {
		t.Fatal(err)
	}
	want = []dbslot.FullStats{
		{UserId: 4, Personal: dbslot.UserStats{SpinsTotal: 1, ScoreCurrent: 5, ScoreBest: 5}, Placement: 1, PlayersTotal: 5},
		{UserId: 2, Personal: dbslot.UserStats{SpinsTotal: 1, ScoreCurrent: 3, ScoreBest: 3}, Placement: 2, PlayersTotal: 5},
		{UserId: 1, Personal: dbslot.UserStats{SpinsTotal: 2, ScoreCurrent: 2, ScoreBest: 2}, Placement: 3, PlayersTotal: 5},
	}
	// New player
	stats, err = dbslot.GetFull(0)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(stats, want) {
		t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
	}
	// Old player in top3
	stats, err = dbslot.GetFull(1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(stats, want) {
		t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
	}
	// Old player outside top3
	want = append(want, dbslot.FullStats{UserId: 5, Personal: dbslot.UserStats{SpinsTotal: 1, ScoreCurrent: 0, ScoreBest: 0}, Placement: 5, PlayersTotal: 5})
	stats, err = dbslot.GetFull(5)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(stats, want) {
		t.Fatalf("unexpected stats:\nwant: %+v\ngot: %+v\n", stats, want)
	}
}
