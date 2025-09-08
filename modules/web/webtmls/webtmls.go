package webtmls

import (
	"html/template"
	"path/filepath"

	"github.com/TrueHopolok/VladOS/modules/cfg"
	"github.com/TrueHopolok/VladOS/modules/db/dbstats"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

// Contain all html templates stored in static/templates/base directory.
// Require [PrepareTemplates] call to be used.
//
// Changing anything about this variable may cause unexpected behaviour if done porly.
var BaseTmls *template.Template

// Reads all base html templates stored in static/tempaltes/base directory.
// Will store the result in [Tmls].
func PrepareBase() error {
	pattern := filepath.Join(cfg.Get().WebStaticDir + "/templates/base/*.html")
	var err error
	BaseTmls, err = template.ParseGlob(pattern)
	return err
}

// Clones existing base templates from [BaseTmls] into a new one.
// Requires for [PrepareBase] to be executed, otherwise they won't be loaded.
//
// Afterwards parses given template names with added prefix of [github.com/TrueHopolok/VladOS/modules/cfg.Cfg.WebStaticPath] + "/templates/".
func ParseTmls(funcMap template.FuncMap, tmlNames ...string) (*template.Template, error) {
	t, err := BaseTmls.Clone()
	if err != nil {
		return t, err
	}
	t.Funcs(funcMap)
	prefixPath := cfg.Get().WebStaticDir + "/templates/"
	for i := range tmlNames {
		tmlNames[i] = prefixPath + tmlNames[i]
	}
	t, err = t.ParseFiles(tmlNames...)
	return t, err
}

// Stores all information that can be used on the page.
type T struct {
	Auth     bool   // must
	Username string // if user is isauth: must; else: optional;
	Title    string // must

	LeaderboardType      string              // if page is leaderboard: must
	LeaderboardPrecent   []dbstats.Precent   // if page is leaderboard: must
	LeaderboardPlacement []dbstats.Placement // if page is leaderboard: must

	SuggestionType string // if page is suggestions: must

	LoginLink string // if page is login: must
	LoginCode string // if page is login: must
}
