package db

import (
	"fmt"
	"io/fs"
)

const MigrationsDir string = "database/migrations/"

// Function checks if any more migrations were added and executes them.
// In case of an empty database, it will create migration table and then perform all available migrations.
//
// Must be used after any initialization function.
func Migrate() error {
	//* Get migration files and migration names
	var fsVersions []string
	// TODO: add file system according to directory
	err := fs.WalkDir("", ".", func(path string, d fs.DirEntry, err error) error {
		return nil // TODO: fill fsVersions to compare later
	})
	if err != nil {
		return fmt.Errorf("walking migration dir error: %w", err)
	}

	//* Prepare db transaction
	tx, err := Conn.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction error: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			panic(fmt.Errorf("transaction rollback error: %w", err))
		}
	}()

	//* Create migration tables if not exists
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS migration (
		version TEXT PRIMARY KEY
	);`)
	if err != nil {
		return fmt.Errorf("create query execution error: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("create query commit error: %w", err)
	}

	//* Recieve all migration versions to compare
	rows, err := tx.Query(`
	SELECT * FROM migration ORDER BY version;
	`)
	if err != nil {
		return fmt.Errorf("select query error: %w", err)
	}
	var dbVersions []string
	for currentVersion := 0; rows.Next(); currentVersion++ {
		dbVersions = append(dbVersions, "")
		err = rows.Scan(&dbVersions[currentVersion])
		if err != nil {
			return fmt.Errorf("scanning rows error: %w", err)
		}
	}

	//* Analyse both migration files and db versions
	if len(dbVersions) > len(fsVersions) {
		return fmt.Errorf("valid migration error: final version of db is unrepeatable since not enough migration files")
	}
	// TODO: finish comparing versions
	// TODO: add executing the migrations logic

	return nil
}
