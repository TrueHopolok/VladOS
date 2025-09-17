package vos_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/TrueHopolok/VladOS/modules/vos"
)

// Helps to check equality of 2 sessions.
func equalSessions(ses1, ses2 vos.Session) bool {
	return ses1.UserID == ses2.UserID && ses1.Username == ses2.Username && ses1.Expire.Equal(ses2.Expire)
}

// WARNING: may fail if fields of the session struct changes.
func TestJson(t *testing.T) {
	ses := vos.NewSession(0, "test", true)
	raw, err := json.Marshal(ses)
	if err != nil {
		t.Fatalf("marshling error: %s", err)
	}
	var got vos.Session
	err = json.Unmarshal(raw, &got)
	if err != nil {
		t.Fatalf("unmarshling error: %s", err)
	}
	if !equalSessions(ses, got) {
		t.Fatalf("corrupeted after marshling")
	}
}

func TestExpiration(t *testing.T) {
	ses := vos.NewSession(0, "test", true)
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
