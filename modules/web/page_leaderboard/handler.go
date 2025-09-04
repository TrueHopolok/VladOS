package page_leaderboard

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

const TmlName string = "leaderboard.html"

var TmlMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
}

func Handle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "leaderboard")

	var (
		data webtmls.T
		ses  vos.Session
		err  error
		tml  *template.Template
	)
	data.Title = "Leaderboard"
	ses, data.Auth = vos.GetSession(r)
	data.Username = ses.Username
	data.Leaderboard, err = dbstats.Leaderboard("dice") // TODO: set to query requested game
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
	}

	tml, err = webtmls.ParseTmls(TmlMap, TmlName)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
	}
	err = tml.ExecuteTemplate(w, TmlName, data)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
	}
}
