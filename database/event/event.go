// /home/krylon/go/src/github.com/blicero/pkman/database/event/event.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-22 20:18:55 krylon>

//go:generate stringer -type=ID

// Package event provides symbolic constants that represent the destructive
// operations we can perform using the package manager.
package event

// ID is the type used to represents events/operations.
type ID uint8

const (
	EventAdd ID = iota
	EventDelete
	EventRefresh
	EventUpdate
	EventClean
	EventAutoremove
)
