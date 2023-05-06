// /home/krylon/go/src/github.com/blicero/pkman/cli/cli.go
// -*- mode: go; coding: utf-8; -*-
// Created on 04. 05. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-05 18:46:46 krylon>

// Package cli implements the command line interface of pkman.
package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/blicero/pkman/backend"
	"github.com/blicero/pkman/common"
	"github.com/blicero/pkman/database"
	"github.com/blicero/pkman/logdomain"
)

// CLI is the nexus of the user interface.
type CLI struct {
	log *log.Logger
	db  *database.Database
	pk  backend.PkgManager
}

// Open creates a new CLI instance.
func Open() (*CLI, error) {
	var (
		err error
		c   = new(CLI)
	)

	if c.log, err = common.GetLogger(logdomain.CLI); err != nil {
		fmt.Fprintf(os.Stdout, "Failed to open Logger for CLI: %s\n",
			err.Error())
		return nil, err
	} else if c.db, err = database.OpenDB(common.DbPath); err != nil {
		c.log.Printf("[ERROR] Cannot open database at %s: %s\n",
			common.DbPath,
			err.Error())
		return nil, err
	}

	return c, nil
} // func Open() (*CLI, error)
