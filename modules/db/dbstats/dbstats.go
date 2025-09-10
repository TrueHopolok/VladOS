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

type Placement struct {
	UserId       int64
	FirstName    string
	Username     string
	Personal     UserStats
	Placement    int
	PlayersTotal int
}

type Precent struct {
	ScoreBest     int
	PlayersAmount int
}

// Return top of all scores and how many players reached that.
func GetPrecent(gameName string) ([]Precent, error) {
	query1, err := QueryDir.ReadFile("precent.sql")
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

	rows, err := tx.Query(fmt.Sprintf(string(query1), "stats_"+gameName))
	if err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return nil, err
	}

	var precent []Precent
	for rows.Next() {
		var p Precent
		if err := rows.Scan(&p.ScoreBest, &p.PlayersAmount); err != nil {
			err = fmt.Errorf("result scanning error: %w", err)
			return nil, err
		}
		precent = append(precent, p)
	}

	return precent, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

// Updates a leaderboard with recieved result for a particular user.
func Update(gameName string, userID int64, firstName string, username string, score int) error {
	query1, err := QueryDir.ReadFile("upd-1.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return err
	}

	query2, err := QueryDir.ReadFile("upd-2.sql")
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

	if _, err := tx.Exec(string(query1), userID, firstName, username); err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return err
	}

	if _, err := tx.Exec(fmt.Sprintf(string(query2), "stats_"+gameName), userID, score); err != nil {
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
func GetSelf(gameName string, userID int64) (UserStats, error) {
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

	rows, err := tx.Query(fmt.Sprintf(string(query), "stats_"+gameName), userID)
	if err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return UserStats{}, err
	}
	if !rows.Next() {
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

func GetTop10(gameName string) ([]Placement, error) {
	return getPlacement(gameName, 0, false)
}

func GetTopSelf(gameName string, userID int64) ([]Placement, error) {
	return getPlacement(gameName, userID, true)
}

func getPlacement(gameName string, userID int64, include bool) ([]Placement, error) {
	var (
		query []byte
		err   error
	)
	if include {
		query, err = QueryDir.ReadFile("placement.sql")
	} else {
		query, err = QueryDir.ReadFile("top10.sql")
	}
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

	var rows *sql.Rows
	if include {
		rows, err = tx.Query(fmt.Sprintf(string(query), "stats_"+gameName), userID)
	} else {
		rows, err = tx.Query(fmt.Sprintf(string(query), "stats_"+gameName))
	}
	if err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return nil, err
	}

	var stats []Placement
	for rows.Next() {
		var next Placement
		if err := rows.Scan(&next.UserId, &next.FirstName, &next.Username, &next.Personal.GamesTotal, &next.Personal.ScoreCurrent, &next.Personal.ScoreBest, &next.Placement, &next.PlayersTotal); err != nil {
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
