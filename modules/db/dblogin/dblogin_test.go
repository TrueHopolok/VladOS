package dblogin_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	"github.com/TrueHopolok/VladOS/modules/db/dblogin"
)

const pathToRoot = "../../../"

type outputData struct {
	UserID    int64
	FirstName string
	Username  string
	Valid     bool
}

func TestLogin(t *testing.T) {
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
		authcode string
		got      outputData
		want     outputData
		err      error
	)

	got.UserID, got.FirstName, got.Username, got.Valid, err = dblogin.Find("anything")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected result\ngot: %v\nwant: %v", got, want)
	}

	authcode, err = dblogin.Add(1, "1", "")
	if err != nil {
		t.Fatal(err)
	}

	got.UserID, got.FirstName, got.Username, got.Valid, err = dblogin.Find("anything")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected result\ngot: %v\nwant: %v", got, want)
	}

	want = outputData{UserID: 1, Valid: true, FirstName: "1", Username: ""}
	got.UserID, got.FirstName, got.Username, got.Valid, err = dblogin.Find(authcode)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected result\ngot: %v\nwant: %v", got, want)
	}

	authcode, err = dblogin.Add(1, "1", "")
	if err != nil {
		t.Fatal(err)
	}

	want = outputData{UserID: 0, Valid: false, FirstName: "", Username: ""}
	got.UserID, got.FirstName, got.Username, got.Valid, err = dblogin.Find("anything")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected result\ngot: %v\nwant: %v", got, want)
	}

	want = outputData{UserID: 1, Valid: true, FirstName: "1", Username: ""}
	got.UserID, got.FirstName, got.Username, got.Valid, err = dblogin.Find(authcode)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected result\ngot: %v\nwant: %v", got, want)
	}

	// wait 5 mins?
	// want = outputData{UserID: 0, Valid: false}
	// got.UserID, got.Valid, err = dblogin.Find(authcode)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if !reflect.DeepEqual(got, want) {
	// 	t.Fatalf("unexpected result\ngot: %v\nwant: %v", got, want)
	// }
}
