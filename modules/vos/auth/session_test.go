package auth_test

import (
	"encoding/json"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/vos/auth"
)

func TestJson(t *testing.T) {
	ses := auth.Session{
		Username: "some name",
	}
	raw, err := json.Marshal(ses)
	if err != nil {
		t.Fatalf("marshling error: %s", err)
	}
	var got auth.Session
	err = json.Unmarshal(raw, &got)
	if err != nil {
		t.Fatalf("unmarshling error: %s", err)
	}
	if ses.Username != got.Username {
		t.Fatalf("corrupeted after marshling")
	}
}

// TODO: add tests for additional functional of session
