// /home/krylon/go/src/github.com/blicero/pkman/database/01_db_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 24. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-24 10:44:57 krylon>

package database

import (
	"testing"

	"github.com/blicero/pkman/common"
)

var db *Database

func TestOpenDB(t *testing.T) {
	var err error

	if db, err = OpenDB(common.DbPath); err != nil {
		db = nil
		t.Fatalf("Error opening database at %s: %s",
			common.DbPath,
			err.Error())
	}
} // func TestOpenDB(t *testing.T)

func TestQueries(t *testing.T) {
	var err error

	if db == nil {
		t.SkipNow()
	}

	for id, str := range qDb {
		if _, err = db.getQuery(id); err != nil {
			t.Errorf("Failed to prepare query %s: %s\n%s",
				id,
				err.Error(),
				str)
		}
	}
} // func TestQueries(t *testing.T)
