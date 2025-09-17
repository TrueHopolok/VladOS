package both_review

import (
	"fmt"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

const TmlPath string = "review/"
const TmlName string = "review_%s.html"
const BaseName string = "review.html"

var TmlMap = template.FuncMap{
	"tn": func(typeName string) string {
		switch typeName {
		case "m8b":
			return "magic 8 ball"
		default:
			return typeName
		}
	},
}

var existingNames = []string{
	"m8b", "tip",
}

func PageHandle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "page/review")

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
	data.Title = "Review"
	ses, data.Auth = vos.GetSession(r)
	data.Username = ses.Username
	data.Admin = ses.Admin
	data.SuggestionType = typeName

	t, err := webtmls.ParseTmls(TmlMap, TmlPath+BaseName, TmlPath+fmt.Sprintf(TmlName, typeName))
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}
	err = t.ExecuteTemplate(w, BaseName, data)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
	}
}
