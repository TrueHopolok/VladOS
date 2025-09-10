package api_auth

import (
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/db/dblogin"
	"github.com/TrueHopolok/VladOS/modules/vos"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "api/login")
	authcode := r.FormValue("authcode")
	userID, firstName, username, valid, err := dblogin.Find(authcode)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}
	if !valid {
		// TODO: add normal error showcase
		http.Error(w, "invalid auth code", http.StatusBadRequest)
		return
	}
	if username == "" {
		username = firstName
	}
	ses := vos.NewSession(userID, username)
	jwt, err := ses.NewJWT()
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}
	vos.SetAuthCookie(w, jwt)
	http.Redirect(w, r, "/", http.StatusFound)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "api/login")
	vos.DeleteAuthCookie(w)
	http.Redirect(w, r, "/", http.StatusFound)
}
