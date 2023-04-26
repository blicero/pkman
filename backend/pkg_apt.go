// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_apt.go
// -*- mode: go; coding: utf-8; -*-
// Created on 21. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-26 11:05:07 krylon>

package backend

import (
	"log"

	"github.com/blicero/krylib"
	"github.com/blicero/pkman/common"
	"github.com/blicero/pkman/database"
	"github.com/blicero/pkman/logdomain"
)

const (
	cmdApt = "/usr/bin/apt"
)

// PkgApt implements the PkgManager interface for Debian's apt.
type PkgApt struct {
	log *log.Logger
	db  *database.Database
}

// CreatePkgApt creates a PkgApt instance to interface with the apt
// package manager used by Debian, Ubuntu, and related distros.
func CreatePkgApt() (*PkgApt, error) {
	var (
		err error
		pk  = new(PkgApt)
	)

	if pk.log, err = common.GetLogger(logdomain.PkgManager); err != nil {
		return nil, err
	} else if pk.db, err = database.OpenDB(common.DbPath); err != nil {
		pk.log.Printf("[ERROR] Cannot open database at %s: %s\n",
			common.DbPath,
			err.Error())
		return nil, err
	}

	return pk, nil
} // func CreatePkgApt() (*PkgApt, error)

func (pk *PkgApt) Search(string) ([]Package, error) {
	return nil, krylib.ErrNotImplemented
} // func (pk *PkgApt) Search(string) ([]Package, error)
