package dbsuggestion

import "embed"

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

//go:embed *.sql
var QueryDir embed.FS

// Add saves provided suggestion from the page
// TODO
func Add(userID int64, typeName string, data []byte) error {
	return nil
}

// Get returns random suggestion to view on the webpage
// TODO
func Get() (any, error) {
	return nil, nil
}
