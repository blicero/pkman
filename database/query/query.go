// /home/krylon/go/src/github.com/blicero/pkman/database/query/query.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-22 19:53:40 krylon>

//go:generate stringer -type=ID

// Package query provides symbolic constants that identify the operations
// we perform on the database.
package query

type ID uint8

const (
	EventAdd ID = iota
	EventGetRecent
	EventGetRecentByType
	EventGetRecentErr
)
