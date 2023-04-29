// /home/krylon/go/src/github.com/blicero/pkman/database/database.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-29 14:27:15 krylon>

// Package database provides the persistence layer and the assorted operations
// we need to perform.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/blicero/krylib"
	"github.com/blicero/pkman/common"
	"github.com/blicero/pkman/database/event"
	"github.com/blicero/pkman/database/query"
	"github.com/blicero/pkman/logdomain"

	_ "github.com/mattn/go-sqlite3" // Import the database driver
)

var (
	openLock sync.Mutex
	retryPat *regexp.Regexp = regexp.MustCompile("(?i)(database is locked|busy)")
)

const (
	retryDelay = 10 * time.Millisecond
)

func worthARetry(err error) bool {
	return retryPat.MatchString(err.Error())
} // func (db *Database) worth_a_retry(err error) bool

func waitForRetry() {
	time.Sleep(retryDelay)
}

// Database wraps the database connection and its associated state and exposes
// the operations we can perform on it.
type Database struct {
	db        *sql.DB
	stmtTable map[query.ID]*sql.Stmt
	tx        *sql.Tx
	log       *log.Logger
	path      string
}

// OpenDB opens a new database connection.
func OpenDB(path string) (*Database, error) {
	var err error
	var msg string
	var dbExists bool

	db := &Database{
		path:      path,
		stmtTable: make(map[query.ID]*sql.Stmt),
	}

	if db.log, err = common.GetLogger(logdomain.Database); err != nil {
		msg = fmt.Sprintf("Error creating logger for Database: %s", err.Error())
		fmt.Println(msg)
		return nil, errors.New(msg)
	}

	var connstring = fmt.Sprintf("%s?_locking=NORMAL&_journal=WAL&_fk=1&recursive_triggers=0",
		path)

	openLock.Lock()
	defer openLock.Unlock()

	if dbExists, err = krylib.Fexists(path); err != nil {
		msg = fmt.Sprintf("Error checking if Database exists at %s: %s", path, err.Error())
		db.log.Println(msg)
		return nil, errors.New(msg)
	} else if db.db, err = sql.Open("sqlite3", connstring); err != nil {
		msg = fmt.Sprintf("Error opening database at %s: %s", path, err.Error())
		db.log.Println(msg)
		return nil, errors.New(msg)
	} else if !dbExists {
		db.log.Printf("Initializing fresh database at %s...\n", path)
		if err = db.initialize(); err != nil {
			msg = fmt.Sprintf("Error initializing database at %s: %s",
				path, err.Error())
			db.log.Println(msg)
			db.db.Close()
			os.Remove(path)
			return nil, errors.New(msg)
		}
	}

	return db, nil
} // func OpenDB(path string) (*Database, error)

func (db *Database) getQuery(qid query.ID) (*sql.Stmt, error) {
	if stmt, ok := db.stmtTable[qid]; ok {
		return stmt, nil
	}

	var stmt *sql.Stmt
	var err error
	var msg string

PREPARE_QUERY:
	if stmt, err = db.db.Prepare(qDb[qid]); err != nil {
		if worthARetry(err) {
			time.Sleep(retryDelay)
			goto PREPARE_QUERY
		} else {
			msg = fmt.Sprintf("Error preparing query %s %s\n\n%s\n",
				qid, err.Error(), qDb[qid])
			db.log.Println(msg)
			return nil, errors.New(msg)
		}
	} else {
		db.stmtTable[qid] = stmt
		return stmt, nil
	}
} // func (db *Database) getQuery(stmt_id query.QueryID) (*sql.Stmt, error)

// Begin starts a transaction
func (db *Database) Begin() error {
	var err error
	var msg string
	var tx *sql.Tx

	if db.tx != nil {
		msg = "Cannot start transaction: A transaction is already in progress!"
		db.log.Println(msg)
		return errors.New(msg)
	}

BEGIN:
	if tx, err = db.db.Begin(); err != nil {
		if worthARetry(err) {
			time.Sleep(retryDelay)
			goto BEGIN
		} else {
			msg = fmt.Sprintf("Cannot start transaction: %s", err.Error())
			db.log.Println(msg)
			return errors.New(msg)
		}
	} else {
		db.tx = tx
		return nil
	}
} // func (db *Database) Begin() error

// Rollback aborts a transaction
func (db *Database) Rollback() error {
	var err error
	var msg string

	if db.tx == nil {
		msg = "Cannot roll back transaction: No transaction is active!"
		db.log.Println(msg)
		return errors.New(msg)
	} else if err = db.tx.Rollback(); err != nil {
		msg = fmt.Sprintf("Cannot roll back transaction: %s", err.Error())
		db.log.Println(msg)
		return errors.New(msg)
	} else {
		db.tx = nil
		return nil
	}
} // func (db *Database) Rollback() error

// Commit finishes a transaction
func (db *Database) Commit() error {
	var err error
	var msg string

	if db.tx == nil {
		msg = "Cannot commit transaction: No transaction is active!"
		db.log.Println(msg)
		return errors.New(msg)
	} else if err = db.tx.Commit(); err != nil {
		msg = fmt.Sprintf("Cannot commit transaction: %s", err.Error())
		db.log.Println(msg)
		return errors.New(msg)
	} else {
		db.tx = nil
		return nil
	}
} // func (db *Database) Commit() error

// Initialize a fresh database, i.e. create all the tables and indices.
// Commit if everythings works as planned, otherwise, roll back, close
// the database, delete the database file, and return an error.
func (db *Database) initialize() error {
	var err error

	err = db.Begin()
	if err != nil {
		msg := fmt.Sprintf("Error starting transaction to initialize database: %s",
			err.Error())
		db.log.Println(msg)
		return errors.New(msg)
	}

	for _, query := range qInit {
		if _, err = db.tx.Exec(query); err != nil {
			msg := fmt.Sprintf("Error executing query %s: %s",
				query, err.Error())
			db.log.Println(msg)
			db.db.Close()
			db.db = nil
			os.Remove(db.path)
			return errors.New(msg)
		}
	}

	db.Commit() // nolint: errcheck
	return nil
} // func (db *Database) initialize() error

// Close closes the database connection
func (db *Database) Close() {
	for _, stmt := range db.stmtTable {
		stmt.Close()
	}

	db.stmtTable = nil

	if db.tx != nil {
		db.tx.Rollback() // nolint: errcheck
		db.tx = nil
	}

	db.db.Close()
} // func (db *Database) Close()

// EventAdd inserts an Event into the database.
func (db *Database) EventAdd(ev *event.Event) error {
	const qid query.ID = query.EventAdd
	var (
		err    error
		msg    string
		stmt   *sql.Stmt
		tx     *sql.Tx
		status bool
	)

	if stmt, err = db.getQuery(qid); err != nil {
		db.log.Printf("[ERROR] Cannot prepare query %s: %s\n",
			qid.String(),
			err.Error())
		return err
	} else if db.tx != nil {
		tx = db.tx
	} else {
	BEGIN_AD_HOC:
		if tx, err = db.db.Begin(); err != nil {
			if worthARetry(err) {
				waitForRetry()
				goto BEGIN_AD_HOC
			} else {
				msg = fmt.Sprintf("Error starting transaction: %s\n",
					err.Error())
				db.log.Printf("[ERROR] %s\n", msg)
				return errors.New(msg)
			}
		} else {
			defer func() {
				var err2 error
				if status {
					if err2 = tx.Commit(); err2 != nil {
						db.log.Printf("[ERROR] Failed to commit ad-hoc transaction: %s\n",
							err2.Error())
					}
				} else if err2 = tx.Rollback(); err2 != nil {
					db.log.Printf("[ERROR] Rollback of ad-hoc transaction failed: %s\n",
						err2.Error())
				}
			}()
		}
	}

	stmt = tx.Stmt(stmt)
	var res sql.Result

EXEC_QUERY:
	if res, err = stmt.Exec(ev.Type, ev.Timestamp.Unix(), ev.Status); err != nil {
		if worthARetry(err) {
			waitForRetry()
			goto EXEC_QUERY
		} else {
			err = fmt.Errorf("Cannot add Event %s to database: %s",
				ev.Type,
				err.Error())
			db.log.Printf("[ERROR] %s\n", err.Error())
			return err
		}
	} else {
		var id int64

		if id, err = res.LastInsertId(); err != nil {
			db.log.Printf("[ERROR] Cannot get ID of new Event %s: %s\n",
				ev.Type,
				err.Error())
			return err
		}

		status = true
		ev.ID = id
		return nil
	}
} // func (db *Database) EventAdd(ev *event.Event) error

// EventGetRecent fetches the (up to) <n> most recent events from the database.
// If n == -1, all Events are fetched.
func (db *Database) EventGetRecent(n int) ([]event.Event, error) {
	const qid query.ID = query.EventGetRecent
	var (
		err  error
		stmt *sql.Stmt
	)

	if stmt, err = db.getQuery(qid); err != nil {
		db.log.Printf("[ERROR] Cannot prepare query %s: %s\n",
			qid,
			err.Error())
		return nil, err
	} else if db.tx != nil {
		stmt = db.tx.Stmt(stmt)
	}

	var rows *sql.Rows

EXEC_QUERY:
	if rows, err = stmt.Query(n); err != nil {
		if worthARetry(err) {
			waitForRetry()
			goto EXEC_QUERY
		}

		return nil, err
	}

	defer rows.Close() // nolint: errcheck,gosec
	var results = make([]event.Event, 0, n)

	for rows.Next() {
		var (
			ev    event.Event
			stamp int64
		)

		if err = rows.Scan(&ev.ID, &ev.Type, &stamp, &ev.Status); err != nil {
			db.log.Printf("[ERROR] Cannot scan row: %s\n", err.Error())
			return nil, err
		}

		ev.Timestamp = time.Unix(stamp, 0)
		results = append(results, ev)
	}

	return results, nil
} // func (db *Database) EventGetRecent(n int) ([]event.Event, error)

// EventGetRecentByType fetches the <n> most recent Events of the given type.
func (db *Database) EventGetRecentByType(n int, evType event.ID) ([]event.Event, error) {
	const qid query.ID = query.EventGetRecentByType
	var (
		err  error
		stmt *sql.Stmt
	)

	if stmt, err = db.getQuery(qid); err != nil {
		db.log.Printf("[ERROR] Cannot prepare query %s: %s\n",
			qid,
			err.Error())
		return nil, err
	} else if db.tx != nil {
		stmt = db.tx.Stmt(stmt)
	}

	var rows *sql.Rows

EXEC_QUERY:
	if rows, err = stmt.Query(n, evType); err != nil {
		if worthARetry(err) {
			waitForRetry()
			goto EXEC_QUERY
		}

		return nil, err
	}

	defer rows.Close() // nolint: errcheck,gosec
	var results = make([]event.Event, 0, n)

	for rows.Next() {
		var (
			ev    = event.Event{Type: evType}
			stamp int64
		)

		if err = rows.Scan(&ev.ID, &stamp, &ev.Status); err != nil {
			db.log.Printf("[ERROR] Cannot scan row: %s\n", err.Error())
			return nil, err
		}

		ev.Timestamp = time.Unix(stamp, 0)
		results = append(results, ev)
	}

	return results, nil
} // func (db *Database) EventGetRecentByType(n int, evType event.ID) ([]event.Event, error)

func (db *Database) EventGetRecentErr(n int) ([]event.Event, error) {
	const qid query.ID = query.EventGetRecentErr
	var (
		err  error
		stmt *sql.Stmt
	)

	if stmt, err = db.getQuery(qid); err != nil {
		db.log.Printf("[ERROR] Cannot prepare query %s: %s\n",
			qid,
			err.Error())
		return nil, err
	} else if db.tx != nil {
		stmt = db.tx.Stmt(stmt)
	}

	var rows *sql.Rows

EXEC_QUERY:
	if rows, err = stmt.Query(n); err != nil {
		if worthARetry(err) {
			waitForRetry()
			goto EXEC_QUERY
		}

		return nil, err
	}

	defer rows.Close() // nolint: errcheck,gosec
	var results = make([]event.Event, 0, n)

	for rows.Next() {
		var (
			ev    event.Event
			stamp int64
		)

		if err = rows.Scan(&ev.ID, &ev.Type, &stamp, &ev.Status); err != nil {
			db.log.Printf("[ERROR] Cannot scan row: %s\n", err.Error())
			return nil, err
		}

		ev.Timestamp = time.Unix(stamp, 0)
		results = append(results, ev)
	}

	return results, nil
} // func (db *Database) EventGetRecentErr(n int) ([]event.Event, error)
