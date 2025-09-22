package dbsuggestion

import (
	"embed"
	"fmt"

	"github.com/TrueHopolok/VladOS/modules/db"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

//go:embed *.sql
var QueryDir embed.FS

// Add saves provided suggestion from the page
func Add(userID int64, typeName string, data []byte) error {
	query, err := QueryDir.ReadFile("add.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		err = fmt.Errorf("beggining connection error: %w", err)
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(string(query), userID, typeName, data)
	if err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return err
	}

	return func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

// Get returns random suggestion to view on the webpage
func GetRandom(typeName string) (id int, userID int64, data string, found bool, err error) {
	query, err := QueryDir.ReadFile("get.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		err = fmt.Errorf("beggining connection error: %w", err)
		return
	}
	defer tx.Rollback()

	rows, err := tx.Query(string(query), typeName)
	if err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return
	}
	found = rows.Next()
	if !found {
		return
	}

	var raw []byte
	if err = rows.Scan(&id, &userID, &raw); err != nil {
		err = fmt.Errorf("result scanning error: %w", err)
		return
	}
	data = string(raw)

	return id, userID, data, found, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

func GetById(typeName string, id int) (string, bool, error) {
	panic("TODO")
	return "", false, nil
}
