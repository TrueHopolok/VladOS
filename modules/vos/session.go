package vos

import "time"

// Contains all necessary information for other packages to use.
// All fields must support [encoding/json] marshalling/unmarshalling and be tested on that.
type Session struct {
	Username string    `json:"Username"`
	Expire   time.Time `json:"Expire"`
}

// Refresh expiration time of the session.
func (ses *Session) Refresh() {
	ses.Expire = time.Now().Add(AuthExpires)
}

// Reports whether or not session is expired.
func (ses Session) Expired() bool {
	return !ses.Expire.After(time.Now())
}

// TODO: add functional for sessions like is expired and etc
