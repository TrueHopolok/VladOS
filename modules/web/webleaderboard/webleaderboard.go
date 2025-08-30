package webleaderboard

import (
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

func Handle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "leaderboard")

	var (
		data webtmls.T
		ses  vos.Session
	)
	data.Title = "Leaderboard"
	ses, data.Auth = vos.GetSession(r)
	data.Username = ses.Username

	err := webtmls.Tmls.ExecuteTemplate(w, "leaderboard.html", data)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
	}
}
