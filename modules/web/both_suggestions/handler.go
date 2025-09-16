package both_suggestions

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/TrueHopolok/VladOS/modules/db/dbsuggestion"
	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

const TmlName string = "suggestions%s.html"

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

	t, err := webtmls.ParseTmls(TmlMap, fmt.Sprintf(TmlName, ""), fmt.Sprintf(TmlName, typeName))
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

	raw := &bytes.Buffer{}
	enc := gob.NewEncoder(raw)
	switch typeName {
	case "tip":
		err := enc.Encode(r.PostFormValue("suggestion"))
		if err != nil {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
			http.Error(w, "http failed", http.StatusInternalServerError)
			return
		}
		err = enc.Encode(r.PostFormValue("author"))
		if err != nil {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
			http.Error(w, "http failed", http.StatusInternalServerError)
			return
		}
	case "m8b":
		err := enc.Encode(r.PostFormValue("suggestion"))
		if err != nil {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
			http.Error(w, "http failed", http.StatusInternalServerError)
			return
		}
		err = enc.Encode(r.PostFormValue("answer"))
		if err != nil {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
			http.Error(w, "http failed", http.StatusInternalServerError)
			return
		}
	default:
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", "unhandled type of suggestion")
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}

	ses, _ := vos.GetSession(r)
	if err := dbsuggestion.Add(ses.UserID, typeName, raw.Bytes()); err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}

	PageHandle(w, r)
}
