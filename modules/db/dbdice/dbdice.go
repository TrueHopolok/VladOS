package dbdice

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/TrueHopolok/VladOS/modules/db"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

//go:embed *.sql
var QueryDir embed.FS

type UserStats struct {
	ThrowsTotal  int
	ScoreCurrent int
	ScoreBest    int
}

// Updates a leaderboard with recieved result for a particular user.
func Update(user_id int64, dice_value int) error {
	query, err := QueryDir.ReadFile("update.sql")
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

	if _, err := tx.Exec(string(query), user_id, dice_value); err != nil {
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

// Recieve stats for certain user and zero if there is no stats.
func Get(user_id int64) (UserStats, error) {
	query, err := QueryDir.ReadFile("get.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return UserStats{}, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		err = fmt.Errorf("beggining connection error: %w", err)
		return UserStats{}, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(string(query), user_id)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserStats{ThrowsTotal: 0, ScoreCurrent: 0, ScoreBest: 0}, nil
		}
		err = fmt.Errorf("query execution error: %w", err)
		return UserStats{}, err
	}
	if !rows.Next() { // should not be possible since there are results and no error, but leave it just in case
		return UserStats{}, nil
	}

	var stats UserStats
	if err := rows.Scan(&stats.ThrowsTotal, &stats.ScoreCurrent, &stats.ScoreBest); err != nil {
		err = fmt.Errorf("result scanning error: %w", err)
		return UserStats{}, err
	}
	return stats, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}
