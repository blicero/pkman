// /home/krylon/go/src/github.com/blicero/pkman/backend/util.go
// -*- mode: go; coding: utf-8; -*-
// Created on 17. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-19 17:31:49 krylon>

package backend

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/blicero/krylib"
)

const releaseFile = "/etc/os-release"

var linePat = regexp.MustCompile(`^(\w+)="([^"]+)"`)

// DetectOS detects the operating system name. Works by calling uname, so this does not
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

// DetectOSVersion returns - if successful - the name and release of the
// operating system we are running on.
func DetectOSVersion() (string, string, error) {
	var (
		err                 error
		line, name, version string
		fh                  *os.File
		rdr                 *bufio.Reader
	)

	if fh, err = os.Open(releaseFile); err != nil {
		return "", "", err
	}

	defer fh.Close() // nolint: errcheck

	rdr = bufio.NewReader(fh)

	for line, err = rdr.ReadString('\n'); err == nil && line != ""; line, err = rdr.ReadString('\n') {
		var match = linePat.FindStringSubmatch(line)

		if len(match) != 3 {
			// fmt.Fprintf(os.Stderr,
			// 	"Could not parse line %q\n",
			// 	line)
			continue
		}

		var key = strings.ToLower(match[1])

		fmt.Fprintf(os.Stderr,
			"%s = %q\n",
			key,
			match[2])

		if key == "name" {
			name = match[2]
		} else if key == "version_id" && version == "" {
			version = match[2]
		} else if key == "version" {
			version = match[2]
		}
	}

	if err == io.EOF {
		err = nil
	}

	return name, version, err
} // func DetectOSVersion() (string, string, error)
