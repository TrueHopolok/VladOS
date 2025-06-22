// This is a sub package of VladOS Security (shortly VOS) package.
//
// Pacakge contains all necessary functional to work with sessions and authefication.
package auth

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

import "time"

// Contains all necessary information for other packages to use.
// All fields must support [encoding/json] marshalling/unmarshalling and be tested on that.
type Session struct {
	Username string    `json:"Username"`
	Expire   time.Time `json:"Expire"`
}
