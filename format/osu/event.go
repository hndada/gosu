package osu

import (
	"errors"
	"strings"
)

// Storyboard is not fully implemented so far.
type Event struct { // delimiter,
	Type      string
	StartTime int
	EndTime   int
	Filename  string
	XOffset   int
	YOffset   int
}

// Exported functions are not guaranteed to be at top of the file.
func newEvent(line string) (ev Event, err error) {
	vs := strings.Split(line, ",")

	switch vs[0] {
	case "0":
		ev.Type = "Background"
	case "1", "Video":
		ev.Type = "Video"
	case "2", "Break":
		ev.Type = "Break"
	}

	switch ev.Type {
	case "Background", "Video":
		if len(vs) < 5 {
			return ev, errors.New("invalid event: not enough length")
		}
		if ev.StartTime, err = parseInt(vs[1]); err != nil {
			return
		}
		ev.Filename = strings.Trim(vs[2], `"`)
		if ev.XOffset, err = parseInt(vs[3]); err != nil {
			return
		}
		if ev.YOffset, err = parseInt(vs[4]); err != nil {
			return
		}

	case "Break":
		if len(vs) < 3 {
			return ev, errors.New("invalid event: not enough length")
		}
		if ev.StartTime, err = parseInt(vs[1]); err != nil {
			return
		}
		if ev.EndTime, err = parseInt(vs[2]); err != nil {
			return
		}
	}

	return
}

func (f Format) Background() (Event, bool) {
	for _, e := range f.Events {
		if e.Type == "Background" {
			return e, true
		}
	}
	return Event{}, false
}

func (f Format) Video() (Event, bool) {
	for _, e := range f.Events {
		if e.Type == "Video" {
			return e, true
		}
	}
	return Event{}, false
}
