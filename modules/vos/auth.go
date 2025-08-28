package vos

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Used to get the session from the request context, but not to block other potential values of the context.
type sessionContextKey struct{}

// Returns [Session] and whether user is autheficated.
// Will always return false in case [AuthMiddleware] was not performed prior.
func GetSession(r *http.Request) (Session, bool) {
	ses := r.Context().Value(sessionContextKey{})
	if ses == nil {
		return Session{}, false
	}
	return ses.(Session), true
}

// After how much time user will be unauthorized.
// Used in both cookies and expiration time in sessions.
const AuthExpires time.Duration = 60 * time.Minute

// Determine which users are allowed to go through [AuthMiddleware] to see given handler.
type AuthFlag int

const (
	// Allow all and any user further to given handler.
	Everyone AuthFlag = 1 << iota

	// Block all unauthorized users from handler.
	Authorized

	// Block authorized users, can be used for login handlers while already authorized.
	Unauthorized
)

func AuthMiddlewareFunc(handler http.HandlerFunc, permissionFlags AuthFlag) http.HandlerFunc {
	return AuthMiddleware(handler, permissionFlags)
}

// Serves as middleware for handlers and users based on authrization permission flags (see [AuthFlag]).
// Will block traffic for unintended users and always update valid auth cookies.
// Logs actions in [log/slog] logger like blocked users or unexpected permissionFlag.
//
// Save authefication status and whole session data in the [net/http.Request.Context].
func AuthMiddleware(handler http.Handler, permissionFlags AuthFlag) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwt, err := GetAuthCookie(r)
		if err != nil {
			DeleteAuthCookie(w)
		}
		ses, isAuthorized, err := ValidateJWT(jwt)
		if err != nil {
			slog.Error(r.Method, "url", r.URL, "auth", "FAILED", "err", err)
			http.Error(w, fmt.Sprintf("failed to get authorization session: %s", err), http.StatusInternalServerError)
			return
		}
		if ses.Expired() {
			isAuthorized = false
		}
		if isAuthorized {
			ses.Refresh()
		}

		switch permissionFlags {
		case Everyone:
			// do nothing since any authorization status is allowed
		case Authorized:
			if !isAuthorized {
				slog.Debug(r.Method, "url", r.URL, "auth", "BLOCKED", "msg", "user is blocked because he is unauthorized")
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
		case Unauthorized:
			if isAuthorized {
				slog.Debug(r.Method, "url", r.URL, "auth", "BLOCKED", "msg", "user is blocked because he is authorized")
				http.Error(w, "user should not be authorized to send this request", http.StatusBadRequest)
				return
			}
		default:
			slog.Error(r.Method, "url", r.URL, "auth", "FAILED", "err", "unknown permission flags, assumed flag AuthFlag.Everyone")
			http.Error(w, "invalid permission flags are selected", http.StatusInternalServerError)
			return
		}

		if isAuthorized {
			r = r.WithContext(context.WithValue(r.Context(), sessionContextKey{}, ses))
			jwt, err = ses.NewJWT()
			if err != nil {
				slog.Error(r.Method, "url", r.URL, "auth", "FAILED", "err", err)
				http.Error(w, fmt.Sprintf("failed to set authorization cookie: %s", err), http.StatusInternalServerError)
				return
			}
			SetAuthCookie(w, jwt)
		}

		slog.Debug(r.Method, "url", r.URL, "auth", "SUCCESS")
		handler.ServeHTTP(w, r)
	}
}
