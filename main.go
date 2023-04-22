// /home/krylon/go/src/github.com/blicero/pkman/main.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-22 14:09:32 krylon>

package main

import (
	"fmt"

	"github.com/blicero/pkman/common"
)

func main() {
	fmt.Printf("%s %s built on %s\n",
		common.AppName,
		common.Version,
		common.BuildStamp.Format(common.TimestampFormat))
}
