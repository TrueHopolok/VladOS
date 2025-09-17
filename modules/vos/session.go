package vos

import "time"

// Contains all necessary information for other packages to use.
// All fields must support [encoding/json] marshalling/unmarshalling and be tested on that.
type Session struct {
	UserID   int64     `json:"UserID"`
	Username string    `json:"Username"`
	Admin    bool      `json:"Admin"`
	Expire   time.Time `json:"Expire"`
}

// Return new session with valid and refreshed expiration time.
func NewSession(userID int64, username string, admin bool) Session {
	return Session{
		UserID:   userID,
		Username: username,
		Admin:    admin,
		Expire:   time.Now().Add(AuthExpires),
	}
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
