// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_pkgadd.go
// -*- mode: go; coding: utf-8; -*-
// Created on 27. 05. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-30 16:10:10 krylon>

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

// PkgOpenBSD implements the PkgManager interface for OpenBSD's binary package
// manager pkg_*
type PkgOpenBSD struct {
	log *log.Logger
	db  *database.Database
}

// CreatePkgOpenBSD creates a new instance of PkgOpenBSD.
func CreatePkgOpenBSD() (*PkgOpenBSD, error) {
	var (
		err error
		pk  = new(PkgOpenBSD)
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
} // func CreatePkgOpenBSD() (*PkgOpenBSD, error)

/* Output of pkg_info -Q emacs:
debug-emacs-28.2p2-gtk2
debug-emacs-28.2p2-gtk3
debug-emacs-28.2p2-no_x11
emacs-28.2p2-gtk2
emacs-28.2p2-gtk3
emacs-28.2p2-no_x11 (installed)
*/

var patSearchPkgOpenBSD = regexp.MustCompile(`(?mi)^(\D+)-(\S+)(?:\s+\(installed\)\s*)?$`)

func (pk *PkgOpenBSD) Search(query string) ([]Package, error) {
	const cmdPkgSearch = "/usr/sbin/pkg_info"
	var (
		err            error
		cmd            *exec.Cmd
		stdout, stderr io.ReadCloser
		bufOut, bufErr bytes.Buffer
	)

	cmd = exec.Command(cmdPkgSearch, "-Q", query)

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
			Name:    m[1],
			Version: m[2],
		}
	}

	return pkList, nil
} // func (pk *PkgOpenBSD) Search(query string) ([]Package, error)

func (pk *PkgOpenBSD) Install(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgOpenBSD) Install(args ...string) error

func (pk *PkgOpenBSD) Remove(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgOpenBSD) Remove(args ...string) error

func (pk *PkgOpenBSD) Update() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgOpenBSD) Update() error

func (pk *PkgOpenBSD) Upgrade() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgOpenBSD) Upgrade() error

func (pk *PkgOpenBSD) ListInstalled() ([]Package, error) {
	return nil, krylib.ErrNotImplemented
} // func (pkg *PkgOpenBSD) ListInstalled() ([]Package, error)

func (pk *PkgOpenBSD) Clean() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgOpenBSD) Clean() error

func (pkg *PkgOpenBSD) LastUpdate() (time.Time, error) {
	return time.Unix(0, 0), krylib.ErrNotImplemented
} // func (pkg *PkgOpenBSD) LastUpdate() (time.Time, error)
