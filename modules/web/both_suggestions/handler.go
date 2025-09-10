package both_suggestions

import (
	"fmt"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

const TmlName string = "suggestions%s.html"

var existingNames = []string{
	"m8b", "tip",
}

func getTypeName(typeName string) string {
	switch typeName {
	case "m8b":
		return "magic 8 ball"
	default:
		return typeName
	}
}

func PageHandle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "page/suggestions")

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

	t, err := webtmls.ParseTmls(template.FuncMap{"tn": getTypeName}, fmt.Sprintf(TmlName, ""), fmt.Sprintf(TmlName, typeName))
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, fmt.Sprintf(TmlName, ""), data)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
	}
}

func PostHandle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "api/suggestions")

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
}
