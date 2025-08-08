package dbslot_test

import (
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbslot"
)

const pathToRoot = "../../../"

func equalStats(x, y dbslot.UserStats) bool {
	return x.SpinsTotal == y.SpinsTotal && x.ScoreCurrent == y.ScoreCurrent && x.ScoreBest == y.ScoreBest
}

func TestSlot(t *testing.T) {
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
