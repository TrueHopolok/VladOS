package dice_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dice"
)

const pathToRoot = "../../../"

func equalStats(x, y dice.UserStats) bool {
	return x.ThrowsTotal == y.ThrowsTotal && x.ThrowsWon == y.ThrowsWon && x.StreakCurrent == y.StreakCurrent && x.StreakBest == y.StreakBest
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

	if err := dice.Update(1, "testname", 1); err != nil {
		t.Fatal(err)
	}
	stats, err := dice.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	want := dice.UserStats{ThrowsTotal: 1, ThrowsWon: 0, StreakCurrent: 0, StreakBest: 0}
	if equalStats(stats, want) {
		t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
	}

	if err := dice.Update(1, "testname", 6); err != nil {
		t.Fatal(err)
	}
	if err := dice.Update(1, "testname", 6); err != nil {
		t.Fatal(err)
	}
	stats, err = dice.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	want = dice.UserStats{ThrowsTotal: 3, ThrowsWon: 2, StreakCurrent: 2, StreakBest: 2}
	if reflect.DeepEqual(stats, want) {
		t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
	}

	if err := dice.Update(1, "testname", 1); err != nil {
		t.Fatal(err)
	}
	if err := dice.Update(1, "testname", 6); err != nil {
		t.Fatal(err)
	}
	stats, err = dice.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	want = dice.UserStats{ThrowsTotal: 5, ThrowsWon: 3, StreakCurrent: 1, StreakBest: 2}
	if reflect.DeepEqual(stats, want) {
		t.Fatalf("Unexpeceted stats:\ngot: %+v\nwant:%+v", stats, want)
	}
}
