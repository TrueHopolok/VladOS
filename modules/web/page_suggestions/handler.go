package page_suggestions

import (
	"log/slog"
	"net/http"

	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

const TmlName string = "suggestions.html"

var existingNames = []string{
	"m8b", "pun", "tip",
}

func Handle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "suggestions")

	typeName := r.URL.Query().Get("type")
	if typeName == "" {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest: ", typeName+" is not recognized")
		http.Error(w, "bad request: no type provided", http.StatusBadRequest)
		return
	}
	found := false
	for _, existingName := range existingNames {
		if typeName == existingName {
			found = true
			break
		}
	}
	if !found {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest: ", typeName+" is not recognized")
		http.Error(w, "bad request: provided type does not supported", http.StatusBadRequest)
		return
	}

	var (
		data webtmls.T
		ses  vos.Session
	)
	data.Title = "Suggestions"
	ses, data.Auth = vos.GetSession(r)
	data.Username = ses.Username
	data.SuggestionType = typeName

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
