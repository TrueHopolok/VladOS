package dbsuggestion_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dbsuggestion"
)

type Result struct {
	Id     int
	UserId int64
	Text   string
}

const pathToRoot = "../../../"

func TestSuggestions(t *testing.T) {
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
		data      []byte
		got, want Result
		found     bool
		err       error
	)

	_, _, _, found, err = dbsuggestion.Rand("1")
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("suggestion in empty table found")
	}

	data, err = json.Marshal("data")
	if err != nil {
		t.Fatal(err)
	}
	if err = dbsuggestion.Add(1, "1", data); err != nil {
		t.Fatal(err)
	}
	_, _, _, found, err = dbsuggestion.Rand("0")
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("suggestion in empty type found")
	}
	got.Id, got.UserId, got.Text, found, err = dbsuggestion.Rand("1")
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("suggestion is not found")
	}
	want = Result{1, 1, `"data"`}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected find 1:\ngot: %v\nwant: %v", got, want)
	}
	_, found, err = dbsuggestion.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("suggestion in empty type found")
	}
	data, found, err = dbsuggestion.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("suggestion is not found")
	}
	got.Text = ""
	if err = json.Unmarshal(data, &got.Text); err != nil {
		t.Fatal(err)
	}
	want.Text = "data"
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected find 2:\ngot: %v\nwant: %v", got, want)
	}

	if err = dbsuggestion.Delete(1); err != nil {
		t.Fatal(err)
	}
	_, found, err = dbsuggestion.Get(1)
	if err != nil {
		t.Fatal(err)
	}
	if found {
		t.Fatal("suggestion in empty table found")
	}
}
