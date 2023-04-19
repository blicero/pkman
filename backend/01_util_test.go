// /home/krylon/go/src/github.com/blicero/pkman/backend/01_util_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 17. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-19 19:11:04 krylon>

package backend

import "testing"

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

func TestDetectOSVersion(t *testing.T) {
	var (
		err           error
		name, version string
	)

	if name, version, err = DetectOSVersion(); err != nil {
		t.Errorf("Failed to detect OS Version: %s",
			err.Error())
	} else if name == "" || version == "" {
		t.Errorf("Failed to detect OS version: %q %q",
			name,
			version)
	} else {
		t.Logf("OS Version is %s %s",
			name,
			version)
	}
} // func TestDetectOSVersion(t *testing.T)
