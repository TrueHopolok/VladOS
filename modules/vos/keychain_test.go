package vos

import (
	"reflect"
	"testing"
)

// Test is inside the package to access all package's data and insure package security/safety by not adding additional testing methods.
//
// Test checks if key chain generates and switches correctly.
func TestChain(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatal("panic", x)
		}
	}()

	var startKeys [2][]byte
	startKeys[0] = make([]byte, EncryptionKeysSize)
	startKeys[1] = make([]byte, EncryptionKeysSize)
	if n := copy(startKeys[0], getCurrentEncryptionKey()); n != EncryptionKeysSize {
		t.Fatalf("start encryption keys size are not equal to packages constant")
	}
	if n := copy(startKeys[1], getPreviousEncryptionKey()); n != EncryptionKeysSize {
		t.Fatalf("start encryption keys size are not equal to packages constant")
	}
	if reflect.DeepEqual(startKeys[0], startKeys[1]) {
		t.Fatalf("start encryption keys are equal, that should not happen")
	}

	switchEncryptionKeys()
	var switchedKeys [2][]byte
	switchedKeys[0] = make([]byte, EncryptionKeysSize)
	switchedKeys[1] = make([]byte, EncryptionKeysSize)
	if n := copy(switchedKeys[0], getCurrentEncryptionKey()); n != EncryptionKeysSize {
		t.Fatalf("switched encryption keys size are not equal to packages constant")
	}
	if n := copy(switchedKeys[1], getPreviousEncryptionKey()); n != EncryptionKeysSize {
		t.Fatalf("switched encryption keys size are not equal to packages constant")
	}
	if reflect.DeepEqual(switchedKeys[0], switchedKeys[1]) {
		t.Fatalf("switched encryption keys are equal, that should not happen")
	}

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
