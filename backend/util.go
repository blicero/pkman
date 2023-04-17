// /home/krylon/go/src/github.com/blicero/pkman/backend/util.go
// -*- mode: go; coding: utf-8; -*-
// Created on 17. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-17 20:43:19 krylon>

package backend

import (
	"os/exec"

	"github.com/blicero/krylib"
)

// Detect the operating system name. Works by calling uname, so this does not
// work on Windows.
func DetectOS() (string, error) {
	var (
		err    error
		name   string
		output []byte
		cmd    = exec.Command("/usr/bin/uname", "-s")
	)

	if output, err = cmd.Output(); err != nil {
		return "", err
	}

	name = string(output)

	return krylib.Chomp(name), nil
} // func DetectOS() (string, error)
