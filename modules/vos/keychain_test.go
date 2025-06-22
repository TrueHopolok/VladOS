package vos_test

import (
	"reflect"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/vos"
)

func TestChain(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()

	startKeys := vos.GetBothEncryptionKeys(t)

	// Test if start chain is valid
	if len(startKeys[0]) != vos.EncryptionKeysSize || len(startKeys[1]) != vos.EncryptionKeysSize {
		t.Fatalf("start encryption keys size are not equal to packages constant")
	}
	if reflect.DeepEqual(startKeys[0], startKeys[1]) {
		t.Fatalf("start encryption keys are equal, that should not happen")
	}

	vos.ManualSwitch(t)
	switchedKeys := vos.GetBothEncryptionKeys(t)

	// Test if switched chain is valid
	if len(switchedKeys[0]) != vos.EncryptionKeysSize || len(switchedKeys[1]) != vos.EncryptionKeysSize {
		t.Fatalf("switched encryption keys size are not equal to packages constant")
	}
	if reflect.DeepEqual(switchedKeys[0], switchedKeys[1]) {
		t.Fatalf("switched encryption keys are equal, that should not happen")
	}

	// Test if switch was performed correctly
	if reflect.DeepEqual(startKeys[0], switchedKeys[0]) {
		t.Fatalf("encryption keys were switched incorrectly, current key was not switched")
	}
	if reflect.DeepEqual(startKeys[1], switchedKeys[1]) {
		t.Fatalf("encryption keys were switched incorrectly, previous key was not switched")
	}
	if !reflect.DeepEqual(startKeys[0], switchedKeys[1]) {
		t.Fatalf("encryption keys were switched incorrectly, start current key didn't become a previous key")
	}
	if reflect.DeepEqual(startKeys[1], switchedKeys[0]) {
		t.Fatalf("encryption keys were switched incorrectly, start previous key somehow became a current key")
	}
}
