// /home/krylon/go/src/github.com/blicero/pkman/backend/01_util_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 17. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-17 21:25:03 krylon>

package backend

import "testing"

func TestDetectOS(t *testing.T) {
	var (
		err error
		s   string
	)

	if s, err = DetectOS(); err != nil {
		t.Errorf("Failed to detect operating system: %s",
			err.Error())
	} else {
		t.Logf("Operating System is %s", s)
	}
} // func TestDetectOS(t *testing.T)
