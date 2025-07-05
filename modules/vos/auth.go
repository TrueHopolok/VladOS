package vos

import (
	"net/http"
	"time"
)

// Determine which users are allowed to go through [AuthMiddleware] to see given handler.
type AuthFlag int

// After how much time user will be unauthorized.
// Used in both cookies and expiration time in sessions.
const AuthExpires time.Duration = 60 * time.Minute

const (
	// Allow all and any user further to given handler.
	Everyone AuthFlag = 1 << iota

	// Block all unauthorized users from handler.
	Authorized

	// Block authorized users, can be used for login handlers while already authorized.
	Unauthorized
)

// Serves as middleware for handlers and users based on authrization permission flags (see [AuthFlag]).
// Will block traffic for unintended users and always update valid auth cookies.
// Logs actions in [log/slog] logger like blocked users or unexpected permissionFlag.
//
// Save authefication status and whole session data in the [net/http.Request.Context].
// TODO: finish this function
func AuthMiddleware(handler http.HandlerFunc, permissionFlags AuthFlag) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwt, err := GetAuthCookie(r)
		if err != nil {
			DeleteAuthCookie(w)
		}
		ses, isAuthorized, err := ValidateJWT(jwt)
		if err != nil {
			// unexpected value
		}

		switch permissionFlags {
		case Everyone:
			// do no checking, since it is allowed for everyone
		case Authorized:
			// block unauthorized traffic
		case Unauthorized:
			// block authorized traffic
		default:
			// assume everyone are allowed, yet warn for unexpected behaviour
		}
		// updateAuthCookie

		handler(w, r)
	}
}
