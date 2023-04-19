// /home/krylon/go/src/github.com/blicero/pkman/backend/platform/system.go
// -*- mode: go; coding: utf-8; -*-
// Created on 19. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-19 22:46:22 krylon>

package platform

// System represents an operating system
type System uint8

// These constants represents the various operating systems we aim to support.
// For systems that have derivatives, such as Debian, Ubuntu, Raspbian, or RHEL
// and Rocky Linux or Alma, we just use one constant, because the package manager
// is going to be the same.
const (
	FreeBSD System = iota
	OpenBSD
	NetBSD
	OpenSuse
	Debian
	Arch
	RedHat
)

// Returns a slice of all System constants.
func AllSystems() []System {
	return []System{
		FreeBSD,
		OpenBSD,
		NetBSD,
		OpenSuse,
		Debian,
		Arch,
		RedHat,
	}
} // func AllSystems() []System
