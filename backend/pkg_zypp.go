// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_zypp.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-21 13:31:16 krylon>

package backend

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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

	fmt.Fprintf(os.Stdout,
		"%s",
		bufOut.String())

	return nil, nil
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
