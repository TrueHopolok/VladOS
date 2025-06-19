// db package used to work in any and all ways with SQLite database.
//
// Provides initalization, migration, execution and testing for database.
package db

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/TrueHopolok/VladOS/modules/cfg"
	_ "github.com/mattn/go-sqlite3"
)

const DBfilePath string = "database/versions/"

type DB struct {
	*sql.DB
}

var Conn DB = DB{nil}

// Connects to SQLite database with path from [github.com/TrueHopolok/VladOS/modules/cfg.Config.DBfilePath].
// Since it is SQLite, the database is opened (or created does not exists) and modified as file.
// Saves in [DB] struct to prevent any outside modifications.
func Init() error {
	conn, err := sql.Open("sqlite3", DBfilePath+cfg.Get().DBfileName)
	if err == nil {
		if conn == nil {
			return fmt.Errorf("database connection is nil")
		}
		Conn = DB{conn}
	} else {
		return err
	}
	return Conn.Ping()
}

// Erase provided database via [os.Remove] and then initialize it like [Init] function.
//
// !WARNING! - must be used only for testing purposes and on testing database to avoid losing data.
func InitTesting(t *testing.T, pathToRoot string) error {
	if !testing.Testing() {
		panic(fmt.Errorf("tried to initialize the database in test mode while not in testing mode"))
	}

	os.Remove(pathToRoot + DBfilePath + cfg.TestGet(pathToRoot).DBfileName)

	conn, err := sql.Open("sqlite3", pathToRoot+DBfilePath+cfg.TestGet(pathToRoot).DBfileName)
	if err == nil {
		if conn == nil {
			return fmt.Errorf("database connection is nil")
		}
		Conn = DB{conn}
	} else {
		return err
	}
	return Conn.Ping()
}
