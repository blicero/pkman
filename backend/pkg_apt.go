// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_apt.go
// -*- mode: go; coding: utf-8; -*-
// Created on 21. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-25 15:34:38 krylon>

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
	cmdApt = "/usr/bin/apt" // nolint: unused
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

/*
Output of apt-cache search emacs (excerpt)
acl2-emacs - Rechenbetonte Logik für applikatives Common Lisp: Emacs-Schnittstelle
xcscope-el - Transition Package, xcscope-el to elpa-xcscope
xemacs21 - highly customizable text editor metapackage
xemacs21-supportel - highly customizable text editor -- non-required library files
xfonts-terminus-oblique - Oblique version of the Terminus font
elpa-xml-rpc - Emacs Lisp XML-RPC client
elpa-xr - convert string regexp to rx notation
elpa-xref - Library for cross-referencing commands in Emacs
xstow - Extended replacement of GNU Stow
elpa-yaml-mode - Emacs major mode for YAML files
elpa-yasnippet - template system for Emacs
yasnippet - transition Package, yasnippet to elpa-yasnippet
elpa-yasnippet-snippets - Andrea Crotti's official YASnippet snippets
yasr - General-purpose console screen reader
yorick - interpreted language and scientific graphics
elpa-zenburn-theme - low contrast color theme for Emacs
elpa-ztree - text mode directory tree
elpa-elfeed-web - Emacs Atom/RSS feed reader - web interface
wnn7egg - Wnn-nana-tamago -- EGG Input Method with Wnn7 for Emacsen
emacs-common-non-dfsg - GNU Emacs common non-DFSG items, including the core documentation
org-mode-doc - keep notes, maintain ToDo lists, and do project planning in emacs

*/

var patSearchApt = regexp.MustCompile(`(?mi)^(\S+)\s+-\s+([^\n]+)$`)

func (pk *PkgApt) Search(query string) ([]Package, error) {
	const cmdSearch = "/usr/bin/apt-cache"
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

	var matches = patSearchApt.FindAllStringSubmatch(bufOut.String(), -1)

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
