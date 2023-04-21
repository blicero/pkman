// /home/krylon/go/src/github.com/blicero/pkman/backend/01_util_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 17. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-21 19:54:08 krylon>

package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseOSRelease(t *testing.T) {
	var (
		err   error
		dirh  *os.File
		files []string
	)

	if dirh, err = os.Open("testdata"); err != nil {
		t.Fatalf("Cannot open testdata folder: %s",
			err.Error())
	}

	defer dirh.Close() // nolint: errcheck

	if files, err = dirh.Readdirnames(-1); err != nil {
		t.Fatalf("Cannot read contents of testdata folder: %s",
			err.Error())
	}

	for _, filename := range files {
		var fpath = filepath.Join("testdata", filename)
		if strings.HasPrefix(filename, "os-release.") {
			fmt.Printf("Attempt to parse %s\n", filename)
			if _, _, err = parseOSRelease(fpath); err != nil {
				t.Errorf("Failed to parse %s: %s",
					filename,
					err.Error())
			}
		}
	}
} // func TestParseOSRelease(t *testing.T)

func TestDetectOS(t *testing.T) {
	var (
		err           error
		name, version string
	)

	if name, version, err = DetectOS(); err != nil {
		t.Errorf("Failed to detect operating system: %s",
			err.Error())
	} else {
		t.Logf("Operating System is %s %s", name, version)
	}
} // func TestDetectOS(t *testing.T)
