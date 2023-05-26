// /home/krylon/go/src/github.com/blicero/pkman/backend/interface.go
// -*- mode: go; coding: utf-8; -*-
// Created on 21. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-26 16:21:21 krylon>

package backend

import (
	"fmt"
	"time"

	"github.com/blicero/pkman/backend/platform"
)

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

// GetPkgManager returns the PkgManager implementation for the given OS.
func GetPkgManager(system string) (PkgManager, error) {
	var (
		err error
		p   platform.System
	)

	if p, err = platform.ParseSystem(system); err != nil {
		return nil, err
	}

	switch p {
	case platform.OpenSuse:
		return CreatePkgZypp()
	case platform.Debian:
		return CreatePkgApt()
	case platform.RedHat:
		return CreatePkgDnf()
	case platform.Arch:
		return CreatePkgPacman()
	default:

		return nil, fmt.Errorf("Support for %s is not implemented", p)
	}
} // func GetPkgManager(system string) (PkgManager, error)
