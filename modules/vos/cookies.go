package vos

import (
	"fmt"
	"net/http"
)

// Name of authefication cookie.
const AuthCookieName string = "auth"

// Read [AuthCookieName] cookie from the request returning its value.
// Expected value is jwt.
// Returns error in case of 0 or >1 cookies were given.
func GetAuthCookie(r *http.Request) (string, error) {
	switch len(r.CookiesNamed(AuthCookieName)) {
	case 0:
		return "", fmt.Errorf("no %s cookies have been received", AuthCookieName)
	case 1:
		// good behaviour, go further in code
	default:
		return "", fmt.Errorf("received too many %s cookies", AuthCookieName)
	}
	cookie, _ := r.Cookie(AuthCookieName)
	return cookie.Value, nil
}

// Sets [AuthCookieName] cookie with given jwt value.
// Expects valid jwt.
// Max age is set to [AuthExpires]
// thus having the same value as default session expiration value.
func SetAuthCookie(w http.ResponseWriter, jwt string) {
	cookie := &http.Cookie{
		Name:     AuthCookieName,
		Value:    jwt,
		Path:     "/",
		MaxAge:   int(AuthExpires.Seconds()),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
}

// Sets [AuthCookieName] cookie max age to -1 thus deleting it from the user.
func DeleteAuthCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     AuthCookieName,
		Value:    "deleted",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
}
