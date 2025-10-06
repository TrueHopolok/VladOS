package websuggestions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"text/template"

	"github.com/TrueHopolok/VladOS/modules/db/dbm8b"
	"github.com/TrueHopolok/VladOS/modules/db/dbpun"
	"github.com/TrueHopolok/VladOS/modules/db/dbsuggestion"
	"github.com/TrueHopolok/VladOS/modules/db/dbtip"
	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

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

var SuggestionExistingNames = []string{
	"m8b", "tip", "pun",
}

// Get suggestion type from request.
//
// If invalid writes into a response body
// After being invalid, further responses are impossible.
func ValidSuggestionType(w http.ResponseWriter, r *http.Request) (valid bool, typeName string) {
	typeName = r.URL.Query().Get("type")
	if typeName == "" {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", typeName+" is not recognized")
		http.Error(w, "bad request: no type provided", http.StatusBadRequest)
		return
	}
	valid = false
	for _, existingName := range SuggestionExistingNames {
		if typeName == existingName {
			valid = true
			break
		}
	}
	if !valid {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", typeName+" is not recognized")
		http.Error(w, "bad request: provided type does not supported", http.StatusBadRequest)
		return
	}
	return
}

func PageHandle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "page/suggestions")

	typeValid, typeName := ValidSuggestionType(w, r)
	if !typeValid {
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

	typeValid, typeName := ValidSuggestionType(w, r)
	if !typeValid {
		return
	}

	var raw []byte
	switch typeName {
	case "tip":
		var sug dbtip.Tip
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
		var sug dbm8b.M8B
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
	case "pun":
		var sug dbpun.Pun
		sug.Pun = r.PostFormValue("pun")
		if len(sug.Pun) == 0 {
			slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", "no pun given")
			http.Error(w, "http failed", http.StatusBadRequest)
			return
		}
		sug.Suffix = r.PostFormValue("suffix")
		if len(sug.Suffix) == 0 {
			slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", "no suffix given")
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
