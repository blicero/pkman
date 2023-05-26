// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_pacman.go
// -*- mode: go; coding: utf-8; -*-
// Created on 25. 05. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-26 16:15:45 krylon>

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

const cmdPacman = "/usr/bin/pacman"

type PkgPacman struct {
	log *log.Logger
	db  *database.Database
}

// CreatePkgPacman creates a PkgPacman instance to interface with the pacman
// package manager used by RHEL, Fedora, and related systems.
func CreatePkgPacman() (*PkgPacman, error) {
	var (
		err error
		pk  = new(PkgPacman)
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
} // func CreatePkgPacman() (*PkgPacman, error)

/*
Output of pacman -Syu emacs

extra/emacs 28.2-2 [Installiert]
    The extensible, customizable, self-documenting real-time display editor
extra/emacs-nativecomp 28.2-2
    The extensible, customizable, self-documenting real-time display editor with native compilation enabled
extra/emacs-nox 28.2-2
    The extensible, customizable, self-documenting real-time display editor without X11 support
community/cl-swank 2.28-1 [Installiert]
    Superior Lisp Interaction Mode for Emacs (Lisp-side server)
community/ecb 2.40.1pre-12
    Emacs Code Browser
community/emacs-apel 10.8.20201107-1
    A library for making portable Emacs Lisp programs.
community/emacs-haskell-mode 17.2-3
    Haskell mode package for Emacs
community/emacs-lua-mode 20210802-3
    Emacs lua-mode
community/emacs-muse 3.20.2-1
    Publishing environment for Emacs
community/emacs-php-mode 1.24.3-1
    PHP mode for emacs
community/emacs-python-mode 6.3.0-1
    Python mode for Emacs
community/emacs-slime 2.28-1 [Installiert]
    Superior Lisp Interaction Mode for Emacs
community/mg 20230406-1
    Micro GNU/emacs
community/semi 1.14.6-9
    A library to provide MIME feature for GNU Emacs.
community/wanderlust 20221010-1
    Mail/News reader supporting IMAP4rev1 for emacs.
*/

var patSearchPacman = regexp.MustCompile(`(?im)^[^/]+/(\S+) ([^\n]+)\s*\n\s+([^\n]+)$`)

func (pk *PkgPacman) Search(query string) ([]Package, error) {
	var (
		err              error
		cmd              *exec.Cmd
		pipeOut, pipeErr io.ReadCloser
		bufOut, bufErr   bytes.Buffer
	)

	cmd = exec.Command(cmdPacman, "-Ss", query)

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

	var matches = patSearchPacman.FindAllStringSubmatch(bufOut.String(), -1)

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
} // func (pk *PkgPacman) Search(string) ([]Package, error)

func (pk *PkgPacman) Install(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgPacman) Install(args ...string) error

func (pk *PkgPacman) Remove(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgPacman) Remove(args ...string) error

func (pk *PkgPacman) Update() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgPacman) Update() error

func (pk *PkgPacman) Upgrade() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgPacman) Upgrade() error

func (pk *PkgPacman) ListInstalled() ([]Package, error) {
	return nil, krylib.ErrNotImplemented
} // func (pkg *PkgPacman) ListInstalled() ([]Package, error)

func (pk *PkgPacman) Clean() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgPacman) Clean() error

func (pkg *PkgPacman) LastUpdate() (time.Time, error) {
	return time.Unix(0, 0), krylib.ErrNotImplemented
} // func (pkg *PkgPacman) LastUpdate() (time.Time, error)
