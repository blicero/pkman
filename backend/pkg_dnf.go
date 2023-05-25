// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_dnf.go
// -*- mode: go; coding: utf-8; -*-
// Created on 25. 05. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-25 16:19:19 krylon>

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

const cmdDnf = "/usr/bin/dnf"

type PkgDnf struct {
	log *log.Logger
	db  *database.Database
}

// CreatePkgDnf creates a PkgDnf instance to interface with the dnf
// package manager used by RHEL, Fedora, and related systems.
func CreatePkgDnf() (*PkgDnf, error) {
	var (
		err error
		pk  = new(PkgDnf)
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
} // func CreatePkgDnf() (*PkgDnf, error)

/*
	Sample output of dnf search:

Last metadata expiration check: 0:00:14 ago on Thu May 25 14:59:46 2023.
======================== Name & Summary Matched: emacs =========================
emacs.x86_64 : GNU Emacs text editor
emacs-auctex.noarch : Enhanced TeX modes for Emacs
emacs-common.x86_64 : Emacs common files
emacs-filesystem.noarch : Emacs filesystem layout
emacs-lucid.x86_64 : GNU Emacs text editor with LUCID toolkit X support
emacs-notmuch.noarch : Not much support for Emacs
emacs-nox.x86_64 : GNU Emacs text editor without X support
poke-emacs.noarch : Emacs support for poke
============================ Summary Matched: emacs ============================
mg.x86_64 : Tiny Emacs-like editor
*/

var patSearchDnf = regexp.MustCompile(`(?im)^(\S+)\s+:\s+([^\n]+)$`)

func (pk *PkgDnf) Search(query string) ([]Package, error) {
	var (
		err              error
		cmd              *exec.Cmd
		pipeOut, pipeErr io.ReadCloser
		bufOut, bufErr   bytes.Buffer
	)

	cmd = exec.Command(cmdDnf, "search", query)

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

	var matches = patSearchDnf.FindAllStringSubmatch(bufOut.String(), -1)

	if len(matches) == 0 {
		return nil, nil
	}

	var pkList = make([]Package, len(matches))

	for i, m := range matches {
		pkList[i] = Package{
			Name:        m[1],
			Description: m[2],
		}
	}

	return pkList, nil
} // func (pk *PkgDnf) Search(string) ([]Package, error)

func (pk *PkgDnf) Install(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgDnf) Install(args ...string) error

func (pk *PkgDnf) Remove(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgDnf) Remove(args ...string) error

func (pk *PkgDnf) Update() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgDnf) Update() error

func (pk *PkgDnf) Upgrade() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgDnf) Upgrade() error

func (pk *PkgDnf) ListInstalled() ([]Package, error) {
	return nil, krylib.ErrNotImplemented
} // func (pkg *PkgDnf) ListInstalled() ([]Package, error)

func (pk *PkgDnf) Clean() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgDnf) Clean() error

func (pkg *PkgDnf) LastUpdate() (time.Time, error) {
	return time.Unix(0, 0), krylib.ErrNotImplemented
} // func (pkg *PkgDnf) LastUpdate() (time.Time, error)
