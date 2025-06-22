package vos

import (
	"crypto/rand"
	"fmt"
	"sync"
	"testing"
	"time"
)

const (
	// Time of a keys switch in minutes.
	EncryptionKeysSwitchingTime time.Duration = 60

	// Size of a signle key in bytes.
	EncryptionKeysSize int = 64
)

// See [KeyChain] docs.
var chain KeyChain

// Contains 2 keys and a mutex, making it multiple goroutines save.
//
// Keys pair consist of:
//   - current key, which is used for an encrypting and decrypting,
//   - previous key, which should be used only for decrypting.
//
// Chain itself and all methods should be used only inside this package.
type KeyChain struct {
	lock sync.Mutex
	keys [2][]byte
}

// Setup the chain value to start values to be ready to use.
// Also start a goroutine to switch encryption keys every [EncryptionKeysSwitchingTime] minutes.
func init() {
	chain.keys[0] = make([]byte, EncryptionKeysSize)
	chain.keys[1] = make([]byte, EncryptionKeysSize)
	rand.Read(chain.keys[0])
	rand.Read(chain.keys[1])
	go cycleSwitch()
}

// Returns the value of a current key in the chain.
//
// This key can be used for both encryption and decryption.
//
// The value of a key should be used only once.
func getCurrentEncryptionKey() []byte {
	chain.lock.Lock()
	defer chain.lock.Unlock()
	return chain.keys[0]
}

// Returns the value of a previous key in the chain.
//
// This key should be used only for old messages decryption.
//
// The value of a key should be used only once.
func getPreviousEncryptionKey() []byte {
	chain.lock.Lock()
	defer chain.lock.Unlock()
	return chain.keys[1]
}

// Automaticly calls [switchEncryptionKeys] every [EncryptionKeysSwitchingTime] minutes
// by repeatadly restarting goroutines.
func cycleSwitch() {
	time.Sleep(EncryptionKeysSwitchingTime * time.Minute)
	go cycleSwitch()
	switchEncryptionKeys()
}

// Switch keys create new current key with
// old one becoming new previous key.
// Previous key will be removed from key chain.
//
// Old instances of keys values should not be used.
func switchEncryptionKeys() {
	chain.lock.Lock()
	defer chain.lock.Unlock()
	copy(chain.keys[1], chain.keys[0])
	rand.Read(chain.keys[0])
}

// Calls [switchEncryptionKeys] outside the time cycle, to test if switch was valid.
//
// !WARNING! - must be used only for testing purposes.
func ManualSwitch(t *testing.T) {
	if !testing.Testing() {
		panic(fmt.Errorf("tried to get raw key chain while not in testing mode"))
	}
	switchEncryptionKeys()
}

// Returns keys field of a encryption key chain to test for validation.
//
// !WARNING! - must be used only for testing purposes.
func GetBothEncryptionKeys(t *testing.T) [2][]byte {
	if !testing.Testing() {
		panic(fmt.Errorf("tried to get raw key chain while not in testing mode"))
	}
	chain.lock.Lock()
	defer chain.lock.Unlock()
	var keysToReturn [2][]byte
	keysToReturn[0] = make([]byte, EncryptionKeysSize)
	keysToReturn[1] = make([]byte, EncryptionKeysSize)
	copy(keysToReturn[0], chain.keys[0])
	copy(keysToReturn[1], chain.keys[1])
	return keysToReturn
}
