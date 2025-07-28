package conversation_test

import (
	"testing"

	"github.com/TrueHopolok/VladOS/modules/db"
	dbconvo "github.com/TrueHopolok/VladOS/modules/db/conversation"
)

const pathToRoot = "../../../"

func TestConversation(t *testing.T) {
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
	if err := dbconvo.Clear(); err != nil {
		t.Fatal(err)
	}

	cs, err := dbconvo.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if !cs.Available {
		t.Fatal("invalid convo status recieved on clear request")
	}

	if err := dbconvo.Free(0); err != nil {
		t.Fatal(err)
	}
	cs, err = dbconvo.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if !cs.Available {
		t.Fatal("invalid convo status recieved on clear+free request")
	}

	if err := dbconvo.Busy(0, "command", nil); err != nil {
		t.Fatal(err)
	}
	cs, err = dbconvo.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if cs.Available {
		t.Fatal("invalid convo status recieved on clear+busy request")
	}

	if err := dbconvo.Free(0); err != nil {
		t.Fatal(err)
	}
	cs, err = dbconvo.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if !cs.Available {
		t.Fatal("invalid convo status recieved on busy+free request")
	}

	if err := dbconvo.Busy(0, "command", nil); err != nil {
		t.Fatal(err)
	}
	if err := dbconvo.Clear(); err != nil {
		t.Fatal(err)
	}
	cs, err = dbconvo.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if !cs.Available {
		t.Fatal("invalid convo status recieved on busy+clear request")
	}
}
