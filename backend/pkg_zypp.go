// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_zypp.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-24 14:40:18 krylon>

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

const cmdZypper = "/usr/bin/zypper"

// PkgZypp implements the PkgManager interface for openSuse's zypper.
type PkgZypp struct {
	log *log.Logger
	db  *database.Database
}

// CreatePkgZypp creates a new instance of PkgZypp.
func CreatePkgZypp() (*PkgZypp, error) {
	var (
		err error
		pk  = new(PkgZypp)
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
} // func CreatePkgZypp() (*PkgZypp, error)

/* Output of zypper search emacs:
Repository 'X11:Utilities' ist veraltet. Sie können 'zypper refresh' als root ausführen, um es zu aktualisieren.
Repository-Daten werden geladen...
Installierte Pakete werden gelesen...

S  | Name                   | Summary                                                        | Type
---+------------------------+----------------------------------------------------------------+------
i+ | emacs                  | GNU Emacs Base Package                                         | Paket
i+ | emacs-apel             | A Portable Emacs Library                                       | Paket
i  | emacs-auctex           | AUC TeX: An Emacs Extension                                    | Paket
   | emacs-color-theme      | Color themes for emacs                                         | Paket
i  | emacs-el               | Several Lisp Files for GNU Emacs                               | Paket
i  | emacs-eln              | GNU Emacs-nox: Emacs Lisp native compiled binary files         | Paket
i+ | emacs-flim             | An Emacs Library for MIME                                      | Paket
   | emacs-gnuplot-mode     | Gnuplot mode for EMACS                                         | Paket
   | emacs-gnuplot-mode-doc | Documentation for EMACS Gnuplot mode                           | Paket
i+ | emacs-info             | Info files for GNU Emacs                                       | Paket
i+ | emacs-nox              | GNU Emacs-nox: An Emacs Binary without X Window System Support | Paket
   | emacs-plugin-devhelp   | Devhelp plugin for Emacs                                       | Paket
   | emacs-poke             | Emacs support for poke                                         | Paket
   | emacs-scheme48         | CMUScheme48 emacs mode                                         | Paket
   | emacs-semi             | Library to provide MIME feature for GNU Emacs                  | Paket
   | emacs-vm               | VM - a mail reader for GNU Emacs                               | Paket
i+ | emacs-w3m              | An interface program to use w3m with Emacs                     | Paket
i+ | emacs-x11              | GNU Emacs: Emacs binary with X Window System Support           | Paket
   | notmuch-emacs          | Emacs lisp email client based on notmuch                       | Paket
   | pinentry-emacs         | Simple PIN or Passphrase Entry Dialog integrated into Emacs    | Paket
   | qemacs                 | An editor similar to Emacs                                     | Paket
   | supercollider-emacs    | SuperCollider support for Emacs                                | Paket
   | vagrant-emacs          | Vagrantfile syntax files for the emacs editor                  | Paket
   | xemacs                 | XEmacs                                                         | Paket
   | xemacs-el              | Emacs-Lisp source files for XEmacs                             | Paket
   | xemacs-info            | Info Files for XEmacs                                          | Paket
   | xemacs-packages        | XEmacs Packages                                                | Paket
   | xemacs-packages-el     | Emacs-Lisp source files for the XEmacs packages                | Paket
   | xemacs-packages-info   | Info Files for the XEmacs Packages                             | Paket

*/

var patSearchZypp = regexp.MustCompile(`(?mi)^(?:i\+?)?\s+\| (\S+)\s+\| (.*?)\s+\| \w+\s*$`)

func (pk *PkgZypp) Search(query string) ([]Package, error) {
	var (
		err            error
		cmd            *exec.Cmd
		stdout, stderr io.ReadCloser
		bufOut, bufErr bytes.Buffer
	)

	cmd = exec.Command(cmdZypper, "se", query)

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
		matches = patSearchZypp.FindAllStringSubmatch(bufOut.String(), -1)
		pkList  = make([]Package, len(matches))
	)

	for i, m := range matches {
		pkList[i] = Package{
			Name:        m[1],
			Description: m[2],
		}
	}

	return pkList, nil
} // func (pk *PkgZypp) Search(query string) ([]Package, error)

func (pk *PkgZypp) Install(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgZypp) Install(args ...string) error

func (pk *PkgZypp) Remove(args ...string) error {
	return krylib.ErrNotImplemented
} // func (pk *PkgZypp) Remove(args ...string) error

func (pk *PkgZypp) Update() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgZypp) Update() error

func (pk *PkgZypp) Upgrade() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgZypp) Upgrade() error

func (pk *PkgZypp) ListInstalled() ([]Package, error) {
	return nil, krylib.ErrNotImplemented
} // func (pkg *PkgZypp) ListInstalled() ([]Package, error)

func (pk *PkgZypp) Clean() error {
	return krylib.ErrNotImplemented
} // func (pk *PkgZypp) Clean() error

func (pkg *PkgZypp) LastUpdate() (time.Time, error) {
	return time.Unix(0, 0), krylib.ErrNotImplemented
} // func (pkg *PkgZypp) LastUpdate() (time.Time, error)
