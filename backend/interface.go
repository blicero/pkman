// /home/krylon/go/src/github.com/blicero/pkman/backend/interface.go
// -*- mode: go; coding: utf-8; -*-
// Created on 21. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-21 23:38:59 krylon>

package backend

import "time"

// PkgManager is a generalized interface to package managers.
type PkgManager interface {
	Search(string) ([]Package, error)
	Install(...string) error
	Remove(...string) error
	Update() error
	Upgrade() error
	ListInstalled() ([]Package, error)
	Clean() error
	LastUpdate() (time.Time, error)
}
