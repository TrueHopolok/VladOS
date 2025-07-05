package vos_test

import (
	"encoding/json"
	"testing"

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

// TODO: add tests for additional functional of session
