package dbpun

import (
	"database/sql"
	"embed"
	"fmt"
	"sync"

	"github.com/TrueHopolok/VladOS/modules/bot/pun/gst"
	"github.com/TrueHopolok/VladOS/modules/db"
)

//go:generate go tool github.com/princjef/gomarkdoc/cmd/gomarkdoc -o documentation.md

//go:embed *.sql
var QueryDir embed.FS

// SyncSuffixTree is just a GST with a [sync.RWMutex].
// Is used for [Write] and [Answer] functions for additional speed.
type SyncSuffixTree struct {
	Mu   sync.RWMutex
	Tree gst.SuffixTree
}

// ArgumentLengthError used by [Write] and [Answer] if any of the given argument is empty.
type ArgumentLengthError struct {
	Description string
}

// Pun used by suggestions to save it into database.
type Pun struct {
	Suffix string `json:"Suffix"`
	Pun    string `json:"Pun"`
}

func (err *ArgumentLengthError) Error() string {
	return err.Description
}

var sst SyncSuffixTree
var once sync.Once

// Write adds a suffix into a database, starting tree and subsequencional gst-s.
func Write(suffix, pun string) error {
	if len(suffix) == 0 {
		return &ArgumentLengthError{"given suffix is empty"}
	}
	if len(pun) == 0 {
		return &ArgumentLengthError{"given pun is empty"}
	}
	if err := loadTree(); err != nil {
		return err
	}

	sst.Mu.Lock()
	defer sst.Mu.Unlock()
	sst.Tree.Put([]byte(suffix))
	raw, err := sst.Tree.Serialize()
	if err != nil {
		return err
	}
	if err = saveTree(raw); err != nil {
		return err
	}
	return addPun(suffix, pun)
}

// Answer returns a random pun choosing the longest suffix of the given word.
func Answer(word string) (pun string, err error) {
	if len(word) == 0 {
		err = &ArgumentLengthError{"given pun is empty"}
		return
	}
	if err = loadTree(); err != nil {
		return
	}

	sst.Mu.RLock()
	suffix := string(sst.Tree.Get([]byte(word)))
	sst.Mu.RUnlock()

	var (
		query []byte
		tx    *sql.Tx
		rows  *sql.Rows
	)
	query, err = QueryDir.ReadFile("get.sql")
	if err != nil {
		err = fmt.Errorf("reading query error: %w", err)
		return
	}

	tx, err = db.Conn.Begin()
	if err != nil {
		err = fmt.Errorf("beggining connection error: %w", err)
		return
	}
	defer tx.Rollback()

	rows, err = tx.Query(string(query), suffix)
	if err != nil {
		err = fmt.Errorf("query execution error: %w", err)
		return
	}
	if !rows.Next() {
		return
	}

	if err = rows.Scan(&pun); err != nil {
		err = fmt.Errorf("result scanning error: %w", err)
		return
	}
	return pun, func() error {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit error: %w", err)
		}
		return nil
	}()
}

// save is writing a raw sst into a database.
// It is expected for raw to be valid serialization of gst.
func saveTree(raw []byte) error {
	query, err := QueryDir.ReadFile("save.sql")
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

	_, err = tx.Exec(string(query), raw)
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

// loadTree executes once by loading GST from database.
func loadTree() (err error) {
	once.Do(func() {
		var (
			query []byte
			tx    *sql.Tx
			rows  *sql.Rows
			raw   []byte
		)
		query, err = QueryDir.ReadFile("init.sql")
		if err != nil {
			err = fmt.Errorf("reading query error: %w", err)
			return
		}

		tx, err := db.Conn.Begin()
		if err != nil {
			err = fmt.Errorf("beggining connection error: %w", err)
			return
		}
		defer tx.Rollback()

		rows, err = tx.Query(string(query))
		if err != nil {
			err = fmt.Errorf("query execution error: %w", err)
			return
		}
		if !rows.Next() {
			return
		}

		if err = rows.Scan(&raw); err != nil {
			err = fmt.Errorf("result scanning error: %w", err)
			return
		}
		if err = tx.Commit(); err != nil {
			err = fmt.Errorf("commit error: %w", err)
			return
		}

		sst.Mu.Lock()
		sst.Tree, err = gst.Deserialize(raw)
		sst.Mu.Unlock()
	})
	return err
}

// addPun adds pun into a db with a key equal to suffix.
func addPun(suffix, pun string) error {
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

	_, err = tx.Exec(string(query), suffix, pun)
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
