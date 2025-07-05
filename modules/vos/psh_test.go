package vos_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/vos"
)

func TestSalt(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()

	salt := vos.GenerateSalt()
	if len(salt) != vos.SaltSize {
		t.Fatal("first salt is not the size of a pacakge's constant")
	}
	if reflect.DeepEqual(salt, vos.GenerateSalt()) {
		t.Fatal("generated salts are equal")
	}
}

func TestPSH(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()

	salt := vos.GenerateSalt()
	password := "Qwerty123"
	psh := vos.NewPSH(password, salt)
	if !vos.ValidatePSH(password, salt, psh) {
		t.Fatalf("generated psh is not valid when should")
	}
	password = "qwerty123"
	if vos.ValidatePSH(password, salt, psh) {
		t.Fatalf("generated psh is valid when should not")
	}
}
