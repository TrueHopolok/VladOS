package dbm8b

import (
	"embed"
	"fmt"

	"github.com/TrueHopolok/VladOS/modules/db"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

type M8B struct {
	Text   string `json:"Text"`
	Answer bool   `json:"Positive"`
}

//go:embed *.sql
var QueryDir embed.FS

// Get random answer depending on recieved result.
func Get(isPositive bool) (string, error) {
	query, err := QueryDir.ReadFile("get.sql")
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

	rows, err := tx.Query(string(query), isPositive)
	if err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return "", err
	}
	if !rows.Next() {
		return "", nil
	}

	var answer string
	if err := rows.Scan(&answer); err != nil {
		err = fmt.Errorf("result scanning error: %w", err)
		return "", err
	}
	return answer, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

func Add(data M8B) error {
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

	_, err = tx.Exec(string(query), data.Answer, data.Text)
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
