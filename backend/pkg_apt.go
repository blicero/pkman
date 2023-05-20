// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_apt.go
// -*- mode: go; coding: utf-8; -*-
// Created on 21. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-20 21:51:10 krylon>

package backend

import (
	"bytes"
	"io"
	"log"
	"os/exec"
	"regexp"
	"time"

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

var patSearch = regexp.MustCompile(`^(\S+) - (.*)`) // nolint: unused

func (pk *PkgApt) Search(query string) ([]Package, error) {
	const cmdSearch = cmdApt // "/usr/bin/apt-cache"
	var (
		err              error
		cmd              *exec.Cmd
		pipeOut, pipeErr io.ReadCloser
		bufOut, bufErr   bytes.Buffer
	)

	cmd = exec.Command(cmdSearch, "search", query)

	if pipeOut, err = cmd.StdoutPipe(); err != nil {
		pk.log.Printf("[ERROR] Cannot get stdout pipe from Cmd: %s\n",
			err.Error())
		return nil, err
	} else if pipeErr, err = cmd.StderrPipe(); err != nil {
		pk.log.Printf("[ERROR] Cannot get stderr pipe from Cmd: %s\n",
			err.Error())
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		pk.log.Printf("[ERROR] Error starting command: %s\n",
			err.Error())
		return nil, err
	}

	io.Copy(&bufOut, pipeOut) // nolint: errcheck
	io.Copy(&bufErr, pipeErr) // nolint: errcheck

	if err = cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			// FIXME Do something with stderr output!
			pk.log.Printf("[ERROR] Failed to wait for command: %s\n",
				err.Error())
			return nil, err
		}
	}

	return nil, krylib.ErrNotImplemented
} // func (pk *PkgApt) Search(string) ([]Package, error)

func (pk *PkgApt) Install(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgApt) Install(args ...string) error

func (pk *PkgApt) Remove(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgApt) Remove(args ...string) error

func (pk *PkgApt) Update() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgApt) Update() error

func (pk *PkgApt) Upgrade() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgApt) Upgrade() error

func (pk *PkgApt) ListInstalled() ([]Package, error) {
	return nil, krylib.ErrNotImplemented
} // func (pkg *PkgApt) ListInstalled() ([]Package, error)

func (pk *PkgApt) Clean() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgApt) Clean() error

func (pkg *PkgApt) LastUpdate() (time.Time, error) {
	return time.Unix(0, 0), krylib.ErrNotImplemented
} // func (pkg *PkgApt) LastUpdate() (time.Time, error)
