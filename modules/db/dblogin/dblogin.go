package dblogin

import (
	"crypto/rand"
	"embed"
	"fmt"
	"time"

	"github.com/TrueHopolok/VladOS/modules/db"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

//go:embed *.sql
var QueryDir embed.FS

// Add to database authcode, clears expired ones and return generated one.
func Add(userID int64, firstName, username string) (string, error) {
	query1, err := QueryDir.ReadFile("add-1-user_data.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return "", err
	}

	query2, err := QueryDir.ReadFile("add-2-drop_old.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return "", err
	}

	query3, err := QueryDir.ReadFile("add-3-find_duplicates.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return "", err
	}

	query4, err := QueryDir.ReadFile("add-4-insert_new.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return "", err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		err = fmt.Errorf("beggining connection error: %w", err)
		return "", err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(query1), userID, firstName, username); err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return "", err
	}

	if _, err := tx.Exec(string(query2), userID, time.Now().Unix()); err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return "", err
	}

	var authcode string
	unique := false
	for !unique {
		authcode = rand.Text()
		rows, err := tx.Query(string(query3), authcode)
		if err != nil {
			err = fmt.Errorf("query execution error: %w", err)
			return "", err
		}
		unique = !rows.Next()
	}

	if _, err := tx.Exec(string(query4), userID, authcode, time.Now().Add(5*time.Minute).Unix()); err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return authcode, err
	}

	return authcode, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}
