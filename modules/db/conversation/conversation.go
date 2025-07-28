package conversation

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/TrueHopolok/VladOS/modules/db"
)

//go:embed *.sql
var QueryDir embed.FS

// Stores info about the user's engagement with the commands.
type Status struct {
	// If user is free from any conversation.
	Available bool

	// Name of command for whom conversation is.
	CommandName string

	// Stores additional data used for conversation like stage of one and previously answered questions.
	Data []byte
}

func Clear() error {
	query, err := QueryDir.ReadFile("clear.sql")
	if err != nil {
		return fmt.Errorf("reading query error: %w", err)
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return fmt.Errorf("beggining connection error: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(query)); err != nil {
		return fmt.Errorf("query execution error: %w", err)
	}
	return func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

func Get(userId int64) (Status, error) {
	res := Status{Available: true, CommandName: "", Data: []byte{}}

	query, err := QueryDir.ReadFile("get.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return res, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		err = fmt.Errorf("beggining connection error: %w", err)
		return res, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(string(query), userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		err = fmt.Errorf("query execution error: %w", err)
		return res, err
	}
	if !rows.Next() { // should not be possible since there are results and no error, but leave it just in case
		return res, nil
	}

	var availableNum int
	if err := rows.Scan(&availableNum, &res.CommandName, &res.Data); err != nil {
		err = fmt.Errorf("result scanning error: %w", err)
		return res, err
	}
	res.Available = availableNum == 1
	return res, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

func Free(userId int64) error {
	query, err := QueryDir.ReadFile("free.sql")
	if err != nil {
		return fmt.Errorf("reading query error: %w", err)
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return fmt.Errorf("beggining connection error: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(query), userId); err != nil {
		return fmt.Errorf("query execution error: %w", err)
	}
	return func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

func Busy(userId int64, cmdName string, data []byte) error {
	query, err := QueryDir.ReadFile("busy.sql")
	if err != nil {
		return fmt.Errorf("reading query error: %w", err)
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		return fmt.Errorf("beggining connection error: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(query), userId, cmdName, data); err != nil {
		return fmt.Errorf("query execution error: %w", err)
	}
	return func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}
