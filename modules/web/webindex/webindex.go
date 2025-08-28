package webindex

import (
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

func Handle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "index-todo")

	var (
		data webtmls.T
		ses  vos.Session
	)
	data.Title = "Home"
	ses, data.Auth = vos.GetSession(r)
	data.Username = ses.Username

	err := webtmls.Tmls.ExecuteTemplate(w, "update.html", data)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
	}
}
