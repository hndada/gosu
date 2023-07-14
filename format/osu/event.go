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
func newEvent(line string) (e Event, err error) {
	vs := strings.Split(line, ",")

	switch vs[0] {
	case "0":
		e.Type = "Background"
	case "1", "Video":
		e.Type = "Video"
	case "2", "Break":
		e.Type = "Break"
	}

	switch e.Type {
	case "Background", "Video":
		if len(vs) < 5 {
			return e, errors.New("invalid event: not enough length")
		}
		if e.StartTime, err = parseInt(vs[1]); err != nil {
			return
		}
		e.Filename = strings.Trim(vs[2], `"`)
		if e.XOffset, err = parseInt(vs[3]); err != nil {
			return
		}
		if e.YOffset, err = parseInt(vs[4]); err != nil {
			return
		}

	case "Break":
		if len(vs) < 3 {
			return e, errors.New("invalid event: not enough length")
		}
		if e.StartTime, err = parseInt(vs[1]); err != nil {
			return
		}
		if e.EndTime, err = parseInt(vs[2]); err != nil {
			return
		}
	}

	return
}

func (es Events) Background() (Event, bool) {
	for _, e := range es {
		if e.Type == "Background" {
			return e, true
		}
	}
	return Event{}, false
}

func (es Events) Video() (Event, bool) {
	for _, e := range es {
		if e.Type == "Video" {
			return e, true
		}
	}
	return Event{}, false
}
