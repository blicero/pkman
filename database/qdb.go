// /home/krylon/go/src/github.com/blicero/pkman/database/qdb.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-22 20:14:36 krylon>

package database

import "github.com/blicero/pkman/database/query"

var qDb = map[query.ID]string{
	query.EventAdd: "INSERT INTO event (event, timestamp, status) VALUES (?, ?, ?)",
	query.EventGetRecent: `
SELECT
    id,
    event,
    timestamp,
    status
FROM event
ORDER BY timestamp DESC
LIMIT ?
`,
	query.EventGetRecentByType: `
SELECT
    id,
    timestamp,
    status
FROM event
WHERE event = ?
ORDER BY timestamp DESC
LIMIT ?
`,
	query.EventGetRecentErr: `
SELECT
    id,
    event,
    timestamp,
    status
FROM event
WHERE status <> 0
ORDER BY timestamp DESC
LIMIT ?
`,
}
