package osu

import (
	"errors"
	"strconv"
	"strings"
)

// storyboard not implemented yet
type Event struct { // delimiter,
	Type      string
	StartTime int
	EndTime   int // optional
	Filename  string
	XOffset   int
	YOffset   int
}

func newEvent(line string) (Event, error) {
	vs := strings.Split(line, `,`)
	var e Event
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
		{
			f, err := strconv.ParseFloat(vs[1], 64)
			if err != nil {
				return e, err
			}
			e.StartTime = int(f)
		}
		{
			e.Filename = strings.Trim(vs[2], `"`)
		}
		{
			f, err := strconv.ParseFloat(vs[3], 64)
			if err != nil {
				return e, err
			}
			e.XOffset = int(f)
		}
		{
			f, err := strconv.ParseFloat(vs[4], 64)
			if err != nil {
				return e, err
			}
			e.YOffset = int(f)
		}
	case "Break":
		{
			if len(vs) < 3 {
				return e, errors.New("invalid event: not enough length")
			}
			f, err := strconv.ParseFloat(vs[1], 64)
			if err != nil {
				return e, err
			}
			e.StartTime = int(f)
		}
		{
			f, err := strconv.ParseFloat(vs[2], 64)
			if err != nil {
				return e, err
			}
			e.EndTime = int(f)
		}
	}
	return e, nil
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
