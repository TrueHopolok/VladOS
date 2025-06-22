package vos

import "time"

// Contains all necessary information for other packages to use.
// All fields must support [encoding/json] marshalling/unmarshalling and be tested on that.
type Session struct {
	Username string    `json:"Username"`
	Expire   time.Time `json:"Expire"`
}
