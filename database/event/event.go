// /home/krylon/go/src/github.com/blicero/pkman/database/event/event.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-22 20:48:03 krylon>

//go:generate stringer -type=ID

// Package event provides symbolic constants that represent the destructive
// operations we can perform using the package manager.
package event

import "time"

// ID is the type used to represents events/operations.
type ID uint8

// These constants represent the operations we can perform on the package manager.
const (
	Add ID = iota
	Delete
	Refresh
	Update
	Clean
	Autoremove
)

// Event represents one operation on the package manager.
type Event struct {
	ID        int64
	Type      ID
	Timestamp time.Time
	Status    int64
}
