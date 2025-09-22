package dbtip

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/TrueHopolok/VladOS/modules/db"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

type Tip struct {
	Text   string `json:"Text"`
	Author string `json:"Author"`
}

//go:embed *.sql
var QueryDir embed.FS

// Retrive a tip from the db table with a given id.
func Get(id int) (text string, author string, found bool, err error) {
	query, err := QueryDir.ReadFile("get.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return "", "", false, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		err = fmt.Errorf("beggining connection error: %w", err)
		return "", "", false, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(string(query), id)
	if err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return "", "", false, err
	}
	if !rows.Next() {
		return "", "", false, nil
	}

	if err := rows.Scan(&text, &author); err != nil {
		err = fmt.Errorf("result scanning error: %w", err)
		return "", "", false, err
	}
	return text, author, true, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

// Retrive a random tip from a db table.
func Rand() (text string, author string, id int, err error) {
	query, err := QueryDir.ReadFile("rand.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return "", "", 0, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		err = fmt.Errorf("beggining connection error: %w", err)
		return "", "", 0, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(string(query))
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", 0, nil
		}
		err = fmt.Errorf("query execution error: %w", err)
		return "", "", 0, err
	}
	if !rows.Next() { // should not be possible since there are results and no error, but leave it just in case
		return "", "", 0, nil
	}

	if err := rows.Scan(&id, &text, &author); err != nil {
		err = fmt.Errorf("result scanning error: %w", err)
		return "", "", 0, err
	}
	return text, author, id, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

func Add(data Tip) error {
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

	_, err = tx.Exec(string(query), data.Author, data.Text)
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
