// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_pkg.go
// -*- mode: go; coding: utf-8; -*-
// Created on 26. 05. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-26 18:39:46 krylon>

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

const cmdPkg = "/usr/sbin/pkg"

// PkgPkg implements the PkgManager interface for openSuse's pkger.
type PkgPkg struct {
	log *log.Logger
	db  *database.Database
}

// CreatePkgPkg creates a new instance of PkgPkg.
func CreatePkgPkg() (*PkgPkg, error) {
	var (
		err error
		pk  = new(PkgPkg)
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
} // func CreatePkgPkg() (*PkgPkg, error)

/* Output of pkg search emacs:
emacs-28.2_4,3                 GNU editing macros
emacs-canna-28.2_4,3           GNU editing macros (Canna Japanese input flavor)
emacs-devel-30.0.50.20230316,3 GNU editing macros
emacs-devel-nox-30.0.50.20230316,3 GNU editing macros (No X flavor)
emacs-koi8u-1.0_1              KOI8-U coding system for [X]Emacs
emacs-lisp-intro-2.04_1        Introduction to Emacs Lisp programming
emacs-nox-28.2_4,3             GNU editing macros (No X flavor)
emacs-w3m-1.4.632.b.20221130   Simple front-end to w3m for emacs
emacs-w3m-emacs_canna-1.4.632.b.20221130 Simple front-end to w3m for emacs
emacs-w3m-emacs_devel-1.4.632.b.20221130 Simple front-end to w3m for emacs
emacs-w3m-emacs_devel_nox-1.4.632.b.20221130 Simple front-end to w3m for emacs
emacs-w3m-emacs_nox-1.4.632.b.20221130 Simple front-end to w3m for emacs
emacsql-3.1.1_2                High-level Emacs Lisp RDBMS front-end
emacsql-emacs_canna-3.1.1_2    High-level Emacs Lisp RDBMS front-end
emacsql-emacs_devel-3.1.1_2    High-level Emacs Lisp RDBMS front-end
emacsql-emacs_devel_nox-3.1.1_2 High-level Emacs Lisp RDBMS front-end
emacsql-emacs_nox-3.1.1_2      High-level Emacs Lisp RDBMS front-end
*/

var patSearchPkg = regexp.MustCompile(`(?mi)^(\S+)-(\d\S+)\s+([^\n]+)$`)

func (pk *PkgPkg) Search(query string) ([]Package, error) {
	var (
		err            error
		cmd            *exec.Cmd
		stdout, stderr io.ReadCloser
		bufOut, bufErr bytes.Buffer
	)

	cmd = exec.Command(cmdPkg, "search", query)

	if stdout, err = cmd.StdoutPipe(); err != nil {
		pk.log.Printf("[ERROR] Cannot get stdout pipe from Cmd: %s\n",
			err.Error())
		return nil, err
	} else if stderr, err = cmd.StderrPipe(); err != nil {
		pk.log.Printf("[ERROR] Cannot get stderr pipe from Cmd: %s\n",
			err.Error())
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		pk.log.Printf("[ERROR] Error starting command: %s\n",
			err.Error())
		return nil, err
	}

	io.Copy(&bufOut, stdout) // nolint: errcheck
	io.Copy(&bufErr, stderr) // nolint: errcheck

	if err = cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			pk.log.Printf("[ERROR] Failed to wait for command: %s\n",
				err.Error())
			return nil, err
		}
	}

	var (
		matches = patSearchPkg.FindAllStringSubmatch(bufOut.String(), -1)
		pkList  = make([]Package, len(matches))
	)

	for i, m := range matches {
		pkList[i] = Package{
			Name:        m[1],
			Version:     m[2],
			Description: m[3],
		}
	}

	return pkList, nil
} // func (pk *PkgPkg) Search(query string) ([]Package, error)

func (pk *PkgPkg) Install(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgPkg) Install(args ...string) error

func (pk *PkgPkg) Remove(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgPkg) Remove(args ...string) error

func (pk *PkgPkg) Update() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgPkg) Update() error

func (pk *PkgPkg) Upgrade() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgPkg) Upgrade() error

func (pk *PkgPkg) ListInstalled() ([]Package, error) {
	return nil, krylib.ErrNotImplemented
} // func (pkg *PkgPkg) ListInstalled() ([]Package, error)

func (pk *PkgPkg) Clean() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgPkg) Clean() error

func (pkg *PkgPkg) LastUpdate() (time.Time, error) {
	return time.Unix(0, 0), krylib.ErrNotImplemented
} // func (pkg *PkgPkg) LastUpdate() (time.Time, error)
