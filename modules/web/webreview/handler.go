package webreview

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/TrueHopolok/VladOS/modules/db/dbm8b"
	"github.com/TrueHopolok/VladOS/modules/db/dbsuggestion"
	"github.com/TrueHopolok/VladOS/modules/db/dbtip"
	"github.com/TrueHopolok/VladOS/modules/vos"
	"github.com/TrueHopolok/VladOS/modules/web/websuggestions"
	"github.com/TrueHopolok/VladOS/modules/web/webtmls"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

const TmlName string = "review.html"

func PostHandle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "post/review")

	typeValid, typeName := websuggestions.ValidSuggestionType(w, r)
	if !typeValid {
		return
	}

	sugTextId := r.URL.Query().Get("id")
	if sugTextId == "" {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", "id is empty")
		http.Error(w, "bad request: no id provided", http.StatusBadRequest)
		return
	}
	suggestionID, err := strconv.Atoi(sugTextId)
	if err != nil {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", err)
		http.Error(w, "bad request: bad id provided", http.StatusBadRequest)
		return
	}

	switch r.FormValue("decision") {
	case "accept":
		// continue with execution
	case "deny":
		if err = dbsuggestion.Delete(suggestionID); err != nil {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
			http.Error(w, "http failed", http.StatusInternalServerError)
			return
		}
		PageHandle(w, r)
	default:
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", "unexpected decision")
		http.Error(w, "bad request: bad id provided", http.StatusBadRequest)
		return
	}

	rawJson, found, err := dbsuggestion.Get(suggestionID)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}
	if !found {
		slog.Debug("http req", "mtd", r.Method, "url", r.URL, "badrequest", "id not found")
		http.Error(w, "bad request: id is not found", http.StatusBadRequest)
		return
	}

	if err = dbsuggestion.Delete(suggestionID); err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}

	switch typeName {
	case "m8b":
		var sug dbm8b.M8B
		if err = json.Unmarshal(rawJson, &sug); err != nil {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", "json unmarshling error")
			http.Error(w, "http failed", http.StatusInternalServerError)
			return
		}
		if err = dbm8b.Add(sug); err != nil {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
			http.Error(w, "http failed", http.StatusInternalServerError)
			return
		}
	case "tip":
		var sug dbtip.Tip
		if err = json.Unmarshal(rawJson, &sug); err != nil {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
			http.Error(w, "http failed", http.StatusInternalServerError)
			return
		}
		if err = dbtip.Add(sug); err != nil {
			slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
			http.Error(w, "http failed", http.StatusInternalServerError)
			return
		}
	default:
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", "unhandled type of suggestion")
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}

	PageHandle(w, r)
}

func PageHandle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("http req", "mtd", r.Method, "url", r.URL, "handler", "page/review")

	typeValid, typeName := websuggestions.ValidSuggestionType(w, r)
	if !typeValid {
		return
	}

	var (
		data webtmls.T
		ses  vos.Session
		err  error
	)
	data.Title = "Review"
	ses, data.Auth = vos.GetSession(r)
	data.Username = ses.Username
	data.Admin = ses.Admin
	data.SuggestionType = typeName
	data.SuggestionID, data.SuggestionUserID, data.SuggestionText, data.SuggestionFound, err = dbsuggestion.Rand(typeName)
	if err != nil {
		slog.Warn("http req", "mtd", r.Method, "url", r.URL, "error", err)
		http.Error(w, "http failed", http.StatusInternalServerError)
		return
	}

	t, err := webtmls.ParseTmls(websuggestions.TmlMap, TmlName)
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
