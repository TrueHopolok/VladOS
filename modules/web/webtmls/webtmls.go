package webtmls

import (
	"html/template"
	"path/filepath"

	"github.com/TrueHopolok/VladOS/modules/cfg"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

// Contain all html templates stored in static/templates directory.
// Require [PrepareTemplates] call to be used.
var Tmls *template.Template

// Reads all html templates stored in static/tempaltes directory.
// Will store the result in [Tmls].
func PrepareTemplates() error {
	pattern := filepath.Join(cfg.Get().WebStaticDir + "/templates/*.html")
	var err error
	Tmls, err = template.ParseGlob(pattern)
	return err
}
