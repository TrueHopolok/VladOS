package vos_test

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/vos"
)

func TestJWT(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()

	ses := vos.Session{
		Username: "test",
	}
	token, err := ses.NewJWT()
	if err != nil {
		t.Fatalf("generating jwt error: %s", err)
	}

	got, valid, err := vos.ValidateJWT(token)
	if err != nil {
		t.Fatalf("1st jwt validation error: %s", err)
	}
	if !valid {
		t.Fatalf("1st jwt validation is incorrect: should have been valid")
	}
	if !reflect.DeepEqual(ses.Username, got.Username) {
		t.Fatalf("1st jwt validation is incorrect: recieved session is not the same")
	}

	jwtHeaderB64 := base64.URLEncoding.EncodeToString([]byte(vos.JWTheader))
	fakeToken := jwtHeaderB64
	sesJson, err := json.Marshal(ses)
	if err != nil {
		t.Fatalf("json-ify session error: %s", err)
	}
	sesB64 := base64.URLEncoding.EncodeToString(sesJson)
	fakeToken += "." + sesB64
	key := make([]byte, vos.EncryptionKeysSize)
	rand.Read(key)
	hash := hmac.New(sha512.New, key)
	hash.Write([]byte(jwtHeaderB64 + sesB64))
	secretB64 := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	fakeToken += "." + secretB64

	_, valid, err = vos.ValidateJWT(fakeToken)
	if err != nil {
		t.Fatalf("2nd jwt validation error: %s", err)
	}
	if valid {
		t.Fatalf("2nd jwt validation is incorrect: should have been invalid")
	}
}
