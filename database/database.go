// /home/krylon/go/src/github.com/blicero/pkman/database/database.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-22 20:24:20 krylon>

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
	"github.com/blicero/pkman/database/query"
)

var (
	openLock sync.Mutex
	retryPat *regexp.Regexp = regexp.MustCompile("(?i)(database is locked|busy)")
)

const (
	retryDelay   = 10 * time.Millisecond
	cacheTimeout = time.Second * 1200
)

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

	if db.log, err = common.GetLogger("Database"); err != nil {
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

func (db *Database) worthARetry(err error) bool {
	return retryPat.MatchString(err.Error())
} // func (db *Database) worth_a_retry(err error) bool

func (db *Database) getStatement(qid query.ID) (*sql.Stmt, error) {
	if stmt, ok := db.stmtTable[qid]; ok {
		return stmt, nil
	}

	var stmt *sql.Stmt
	var err error
	var msg string

PREPARE_QUERY:
	if stmt, err = db.db.Prepare(qDb[qid]); err != nil {
		if db.worthARetry(err) {
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
} // func (db *Database) getStatement(stmt_id query.QueryID) (*sql.Stmt, error)

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
		if db.worthARetry(err) {
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
