package vos

import (
	"crypto/rand"
	"sync"
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
// Used by JWT and a few additional encryption functions.
//
// Keys pair switch keys every [EncryptionKeysSwitchingTime] minutes, thus providing additional security.
//
// Keys pair consist of:
//   - current key, which is used for an encrypting and decrypting,
//   - previous key, which should be used only for decrypting.
//
// Chain itself and all its methods should be used only inside VOS package.
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
	go func() {
		for {
			time.Sleep(EncryptionKeysSwitchingTime * time.Minute)
			switchEncryptionKeys()
		}
	}()
}

// Returns the copied value of a current key in the chain.
//
// This key can be used for both encryption and decryption.
//
// The value of a key should be used only once.
func getCurrentEncryptionKey() []byte {
	chain.lock.Lock()
	defer chain.lock.Unlock()
	key := make([]byte, EncryptionKeysSize)
	copy(key, chain.keys[0])
	return key
}

// Returns the copied value of a previous key in the chain.
//
// This key should be used only for old messages decryption.
//
// The value of a key should be used only once.
func getPreviousEncryptionKey() []byte {
	chain.lock.Lock()
	defer chain.lock.Unlock()
	key := make([]byte, EncryptionKeysSize)
	copy(key, chain.keys[1])
	return key
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
