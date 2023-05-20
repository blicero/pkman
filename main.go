// /home/krylon/go/src/github.com/blicero/pkman/main.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-05-20 12:30:57 krylon>

package main

import (
	"fmt"
	"os"

	"github.com/blicero/pkman/cli"
	"github.com/blicero/pkman/common"
)

func main() {
	fmt.Printf("%s %s built on %s\n",
		common.AppName,
		common.Version,
		common.BuildStamp.Format(common.TimestampFormat))
	var (
		err error
		c   *cli.CLI
	)

	if c, err = cli.Open(); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Cannot open CLI: %s\n",
			err.Error())
		os.Exit(1)
	}

	c.Run()
}
