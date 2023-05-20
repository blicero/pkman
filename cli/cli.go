// /home/krylon/go/src/github.com/blicero/pkman/cli/cli.go
// -*- mode: go; coding: utf-8; -*-
// Created on 04. 05. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-20 21:46:12 krylon>

// Package cli implements the command line interface of pkman.
package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

func (c *CLI) Run() {
	var (
		err               error
		name, release, op string
		args              []string
		pk                backend.PkgManager
	)

	if name, release, err = backend.DetectOS(); err != nil {
		c.log.Printf("[ERROR] Cannot detect operating system: %s\n",
			err.Error())
		return
	} else if pk, err = backend.GetPkgManager(name); err != nil {
		c.log.Printf("[ERROR] Failed to get PkgManager for %s: %s\n",
			name,
			err.Error())
		return
	}

	c.log.Printf("[DEBUG] We are running on %s %s\n",
		name,
		release)

	flag.Parse()

	args = flag.Args()

	if len(args) == 0 {
		fmt.Println("No arguments were given. So long.")
		os.Exit(0)
	}

	op = args[0]
	args = args[1:]

	switch strings.ToLower(op) {
	case "se", "search":
		var pkList []backend.Package

		if pkList, err = pk.Search(args[0]); err != nil {
			c.log.Printf("[ERROR] Failed to search for %q: %s\n",
				args[0],
				err.Error())
		} else {
			for _, p := range pkList {
				fmt.Printf("%20s %10s %s\n",
					p.Name,
					p.Version,
					p.Description)
			}
		}
	default:
		c.log.Printf("[ERROR] Unsupported operation %q\n",
			op)
	}
}
