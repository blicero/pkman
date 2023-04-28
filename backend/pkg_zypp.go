// /home/krylon/go/src/github.com/blicero/pkman/backend/pkg_zypp.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-28 11:45:24 krylon>

package backend

import (
	"bytes"
	"io"
	"log"
	"os/exec"

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

	cmd = exec.Command(cmdApt, "se", query)

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

	return nil, krylib.ErrNotImplemented
} // func (pk *PkgZypp) Search(query string) ([]Package, error)
