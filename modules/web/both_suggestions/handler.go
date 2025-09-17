package both_suggestions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/TrueHopolok/VladOS/modules/db/dbsuggestion"
	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

type SuggestionTip struct {
	Text   string `json:"Text"`
	Author string `json:"Author"`
}

type SuggestionM8B struct {
	Text   string `json:"Text"`
	Answer bool   `json:"Positive"`
}

const TmlPath string = "suggestions/"
const TmlName string = "suggest_%s.html"
const BaseName string = "suggest.html"

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

func PostHandle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "api/suggestions")

	typeName := r.URL.Query().Get("type")
	if typeName == "" {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", typeName+" is not recognized")
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
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", typeName+" is not recognized")
		http.Error(w, "bad request: provided type does not supported", http.StatusBadRequest)
		return
	}

	var raw []byte
	switch typeName {
	case "tip":
		var sug SuggestionTip
		sug.Text = r.PostFormValue("suggestion")
		if len(sug.Text) == 0 {
			slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", "no suggestion given")
			http.Error(w, "http failed", http.StatusBadRequest)
			return
		}
		sug.Author = r.PostFormValue("author")
		var err error
		raw, err = json.MarshalIndent(sug, "", "  ")
		if err != nil {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
			http.Error(w, "http failed", http.StatusInternalServerError)
			return
		}
	case "m8b":
		var sug SuggestionM8B
		sug.Text = r.PostFormValue("suggestion")
		if len(sug.Text) == 0 {
			slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", "no suggestion given")
			http.Error(w, "http failed", http.StatusBadRequest)
			return
		}
		switch r.PostFormValue("answer") {
		case "yes":
			sug.Answer = true
		case "no":
			sug.Answer = false
		default:
			slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", "invalid answer is given")
			http.Error(w, "http failed", http.StatusBadRequest)
			return
		}
		var err error
		raw, err = json.MarshalIndent(sug, "", "  ")
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
	bts := &bytes.Buffer{}
	json.HTMLEscape(bts, raw)

	ses, _ := vos.GetSession(r)
	if err := dbsuggestion.Add(ses.UserID, typeName, bts.Bytes()); err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}

	PageHandle(w, r)
}
