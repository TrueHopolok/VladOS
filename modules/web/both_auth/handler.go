package both_auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/db/dblogin"
	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

type errorContextKey struct{}

const TmlName string = "login.html"

func PageHandle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "page/login")

	var (
		data webtmls.T
	)
	data.Title = "Login"
	data.Auth = false
	errText := r.Context().Value(errorContextKey{})
	if errText != nil {
		data.LoginError = errText.(string)
	}
	fmt.Println(data.LoginError)

	t, err := webtmls.ParseTmls(nil, TmlName)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, TmlName, data)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
	}
}

func LoginHandle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "api/login")
	authcode := r.FormValue("authcode")
	userID, firstName, username, valid, err := dblogin.Find(authcode)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}
	if !valid {
		r = r.WithContext(context.WithValue(r.Context(), errorContextKey{}, "invalid auth code"))
		PageHandle(w, r)
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

func LogoutHandle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "api/login")
	vos.DeleteAuthCookie(w)
	http.Redirect(w, r, "/", http.StatusFound)
}
