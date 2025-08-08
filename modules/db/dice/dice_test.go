package dice_test

import (
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dice"
)

const pathToRoot = "../../../"

func equalStats(x, y dice.UserStats) bool {
	return x.ThrowsTotal == y.ThrowsTotal && x.ScoreCurrent == y.ScoreCurrent && x.ScoreBest == y.ScoreBest
}

func TestDice(t *testing.T) {
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

	if err := dice.Update(1, 1); err != nil {
		t.Fatal(err)
	}
	stats, err := dice.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	want := dice.UserStats{ThrowsTotal: 1, ScoreCurrent: 0, ScoreBest: 0}
	if !equalStats(stats, want) {
		t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
	}

	if err := dice.Update(1, 3); err != nil {
		t.Fatal(err)
	}
	if err := dice.Update(1, 6); err != nil {
		t.Fatal(err)
	}
	stats, err = dice.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	want = dice.UserStats{ThrowsTotal: 3, ScoreCurrent: 9, ScoreBest: 9}
	if !equalStats(stats, want) {
		t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
	}

	if err := dice.Update(1, 1); err != nil {
		t.Fatal(err)
	}
	if err := dice.Update(1, 6); err != nil {
		t.Fatal(err)
	}
	stats, err = dice.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	want = dice.UserStats{ThrowsTotal: 5, ScoreCurrent: 6, ScoreBest: 9}
	if !equalStats(stats, want) {
		t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
	}
}
