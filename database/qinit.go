// /home/krylon/go/src/github.com/blicero/pkman/database/qinit.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-29 14:26:44 krylon>

package database

var qInit = []string{
	`
CREATE TABLE event (
    id		INTEGER PRIMARY KEY,
    event	INTEGER NOT NULL,
    timestamp	INTEGER NOT NULL,
    status	INTEGER NOT NULL
) STRICT
`,
	"CREATE INDEX ev_event_idx ON event (event)",
	"CREATE INDEX ev_timestamp_idx ON event (timestamp)",
	"CREATE INDEX ev_status_idx ON event (status)",
}
