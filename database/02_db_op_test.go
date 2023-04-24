// /home/krylon/go/src/github.com/blicero/pkman/database/02_db_op_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 24. 04. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-04-24 20:25:04 krylon>

package database

import (
	"math/rand"
	"testing"
	"time"

	"github.com/blicero/pkman/database/event"
)

const evCnt = 168

var (
	initEvents [evCnt]event.Event
	evStamp    = time.Now().Add(time.Hour * -24 * 7)
)

func randomEvent(ev *event.Event) {
	ev.Type = event.ID(rand.Intn(event.EventCnt))
	ev.Timestamp = evStamp
	// ev.Status =
	if rand.Intn(100) < 75 {
		ev.Status = 0
	} else {
		ev.Status = rand.Int63n(254) + 1
	}
	evStamp = evStamp.Add(time.Hour)
} // func randomEvent(ev *event.Event)

func TestEventAdd(t *testing.T) {
	if db == nil {
		t.SkipNow()
	}

	for i := range initEvents {
		var (
			err error
			ev  = &initEvents[i]
		)

		randomEvent(ev)

		if err = db.EventAdd(ev); err != nil {
			t.Errorf("Failed to add Event #%d: %s",
				i,
				err.Error())
		} else if ev.ID == 0 {
			t.Error("EventAdd did not return an error, but it did not set an ID")
		}
	}
} // func TestEventAdd(t *testing.T)

func TestEventGetRecent(t *testing.T) {
	if db == nil {
		t.SkipNow()
	}

	for i := 0; i < 20; i++ {
		var (
			err    error
			evList []event.Event
			cnt    = rand.Intn(evCnt) + 1
		)

		if evList, err = db.EventGetRecent(cnt); err != nil {
			t.Errorf("Error fetching %d recent events: %s",
				cnt,
				err.Error())
		} else if len(evList) != cnt {
			t.Errorf("Unexpected number of results for EventGetRecent: %d (expected %d)",
				len(evList),
				cnt)
		}
	}
} // func TestEventGetRecent(t *testing.T)
