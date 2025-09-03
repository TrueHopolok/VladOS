package page_index

import (
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

const TmlName string = "update.html"

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "index")

	var (
		data webtmls.T
		ses  vos.Session
	)
	data.Title = "Home"
	ses, data.Auth = vos.GetSession(r)
	data.Username = ses.Username

	t, err := webtmls.ParseTmls(nil, TmlName)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
	}
	err = t.ExecuteTemplate(w, TmlName, data)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
	}
}
