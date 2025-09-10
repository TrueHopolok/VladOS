package page_login

import (
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

const TmlName string = "login.html"

func Handle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "page/login")

	var (
		data webtmls.T
	)
	data.Title = "Login"
	data.Auth = false

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
