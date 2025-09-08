package dbstats_test

const pathToRoot = "../../../"

var tablesToTest = []string{"slot", "dice", "bjack", "guess"}

/*
func TestExample(t *testing.T) {
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
		want = nil
		got, err = dbstats.GetTEST(gameName)
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
		got, err = dbstats.GetTEST(gameName)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected result\ngot: %v\nwant:%v", got, want)
		}
	}
}
*/
