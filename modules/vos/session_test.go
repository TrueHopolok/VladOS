package vos_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TrueHopolok/VladOS/modules/vos"
)

func TestJson(t *testing.T) {
	ses := vos.Session{
		Username: "some name",
	}
	raw, err := json.Marshal(ses)
	if err != nil {
		t.Fatalf("marshling error: %s", err)
	}
	var got vos.Session
	err = json.Unmarshal(raw, &got)
	if err != nil {
		t.Fatalf("unmarshling error: %s", err)
	}
	if ses.Username != got.Username {
		t.Fatalf("corrupeted after marshling")
	}
}

func TestExpiration(t *testing.T) {
	ses := vos.Session{
		Expire: time.Now().Add(vos.AuthExpires),
	}
	if ses.Expired() {
		t.Fatalf("new seesion expired when should not")
	}
	ses.Expire = time.Now().Add(-1)
	if !ses.Expired() {
		t.Fatalf("old session was not expired")
	}
	ses.Refresh()
	if ses.Expired() {
		t.Fatalf("refreshed seesion expired when should not")
	}
}

// TODO: add tests for additional functional of session
