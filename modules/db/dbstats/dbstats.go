package dbstats

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
	GamesTotal   int
	ScoreCurrent int
	ScoreBest    int
}

type FullStats struct {
	UserId       int64
	Personal     UserStats
	Placement    int
	PlayersTotal int
}

type Placement struct {
	ScoreBest     int
	PlayersAmount int
}

// Return top of all scores and how many players reached that.
func Leaderboard(gameName string) ([]Placement, error) {
	query, err := QueryDir.ReadFile("leaderboard.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return nil, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		err = fmt.Errorf("beggining connection error: %w", err)
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(fmt.Sprintf(string(query), gameName))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		err = fmt.Errorf("query execution error: %w", err)
		return nil, err
	}

	var leaderboard []Placement
	for rows.Next() {
		var p Placement
		if err := rows.Scan(&p.ScoreBest, &p.PlayersAmount); err != nil {
			err = fmt.Errorf("result scanning error: %w", err)
			return nil, err
		}
		leaderboard = append(leaderboard, p)
	}

	return leaderboard, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

// Updates a leaderboard with recieved result for a particular user.
func Update(gameName string, userID int64, score int) error {
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

	if _, err := tx.Exec(fmt.Sprintf(string(query), gameName), userID, score); err != nil {
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
func Get(gameName string, userID int64) (UserStats, error) {
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

	rows, err := tx.Query(fmt.Sprintf(string(query), gameName), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserStats{}, nil
		}
		err = fmt.Errorf("query execution error: %w", err)
		return UserStats{}, err
	}
	if !rows.Next() { // should not be possible since there are results and no error, but leave it just in case
		return UserStats{}, nil
	}

	var stats UserStats
	if err := rows.Scan(&stats.GamesTotal, &stats.ScoreCurrent, &stats.ScoreBest); err != nil {
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

func GetFull(gameName string, userID int64) ([]FullStats, error) {
	query, err := QueryDir.ReadFile("full.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return nil, err
	}

	tx, err := db.Conn.Begin()
	if err != nil {
		err = fmt.Errorf("beggining connection error: %w", err)
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(fmt.Sprintf(string(query), gameName), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		err = fmt.Errorf("query execution error: %w", err)
		return nil, err
	}

	var stats []FullStats
	for rows.Next() {
		var next FullStats
		if err := rows.Scan(&next.UserId, &next.Personal.GamesTotal, &next.Personal.ScoreCurrent, &next.Personal.ScoreBest, &next.Placement, &next.PlayersTotal); err != nil {
			err = fmt.Errorf("result scanning error: %w", err)
			return nil, err
		}
		stats = append(stats, next)
	}
	return stats, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}
