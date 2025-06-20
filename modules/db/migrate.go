package db

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
)

//go:embed migrations/*.sql
var MigrationsDir embed.FS

// Function checks if any more migrations were added and executes them.
// In case of an empty database, it will create migration table and then perform all available migrations.
//
// Must be used after any initialization function.
func Migrate() error {
	//* Get migration files and migration names
	var fsVersions []string
	err := fs.WalkDir(MigrationsDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fsVersions = append(fsVersions, path)
			slog.Debug("db migrate", "FOUND migration file", path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walking embed migration dir error: %w", err)
	}
	slog.Debug("db migrate", "AMOUNT of found migrations files", len(fsVersions))

	//* Create migration tables if not exists
	_, err = Conn.Exec(`
	CREATE TABLE IF NOT EXISTS migration (
		version TEXT PRIMARY KEY
	);`)
	if err != nil {
		return fmt.Errorf("create query execution error: %w", err)
	}

	//* Recieve all migration versions to compare
	rows, err := Conn.Query(`
	SELECT * FROM migration ORDER BY version;
	`)
	if err != nil {
		return fmt.Errorf("select query error: %w", err)
	}
	var dbVersions []string
	for i := 0; rows.Next(); i++ {
		dbVersions = append(dbVersions, "")
		err = rows.Scan(&dbVersions[i])
		if err != nil {
			return fmt.Errorf("scanning rows error: %w", err)
		}
		slog.Debug("db migrate", "FOUND db version", dbVersions[i])
	}
	slog.Debug("db migrate", "AMOUNT of found db versions", len(dbVersions))

	//* Analyse both migration files and db versions
	if len(dbVersions) > len(fsVersions) {
		return fmt.Errorf("final version of db is unrepeatable: not enough migration files")
	}
	i := 0
	for ; i < len(dbVersions); i++ {
		if dbVersions[i] != fsVersions[i] {
			return fmt.Errorf("final version of db is unrepetable: DB migrations history differs from FS history: DB=%s FS=%s", dbVersions[i], fsVersions[i])
		}
	}

	//* Execute missing migrations
	for ; i < len(fsVersions); i++ {
		migrationVersion := fsVersions[i]
		slog.Debug("db migrate", "EXECUTE migration", migrationVersion, "STATUS", "START")
		rawMigration, err := MigrationsDir.ReadFile(migrationVersion)
		if err != nil {
			return fmt.Errorf("migration %s READING error: %w", migrationVersion, err)
		}
		if err = runMigration(migrationVersion, rawMigration); err != nil {
			return fmt.Errorf("migration %s EXECUTION error: %w", migrationVersion, err)
		}
		slog.Debug("db migrate", "EXECUTE migration", migrationVersion, "STATUS", "SUCCESS")
	}

	return nil
}

// Execute single migration given as a version and query itself.
// Will update both the database with that migration and migration table itself.
func runMigration(migrationVersion string, rawMigration []byte) error {
	tx, err := Conn.Begin()
	if err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(string(rawMigration)); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	if _, err = tx.Exec(fmt.Sprintf(`
	INSERT INTO migration (version) VALUES ("%s");
	`, migrationVersion)); err != nil {
		return fmt.Errorf("version update error: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}
