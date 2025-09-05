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

var existingNames = []string{
	"dice", "guess", "bjack", "slot",
}

func Handle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "leaderboard")

	gameName := r.URL.Query().Get("game")
	if gameName == "" {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest: ", gameName+" is not recognized")
		http.Error(w, "bad request: no game provided", http.StatusBadRequest)
		return
	}
	found := false
	for _, existingName := range existingNames {
		if gameName == existingName {
			found = true
			break
		}
	}
	if !found {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest: ", gameName+" is not recognized")
		http.Error(w, "bad request: provided game does not supported", http.StatusBadRequest)
		return
	}

	var (
		data webtmls.T
		ses  vos.Session
		err  error
		tml  *template.Template
	)
	data.Title = "Leaderboard"
	ses, data.Auth = vos.GetSession(r)
	data.Username = ses.Username
	data.LeaderboardType = gameName
	data.Leaderboard, err = dbstats.Leaderboard(gameName)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}

	tml, err = webtmls.ParseTmls(TmlMap, TmlName)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}
	err = tml.ExecuteTemplate(w, TmlName, data)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}
}
