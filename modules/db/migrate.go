package db

import "os"

const MigrationsDir string = "database/migrations/"

func Migrate() error {
	// TODO: finish migrations
	_, err := os.ReadDir(MigrationsDir)
	return err
}
