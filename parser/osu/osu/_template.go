package osu

import (
	"bytes"
	"io/ioutil"
	"strings"
)

// todo: put flag: force parse (ignore error)
func Parse(path string) (*FormatOsu, error) {
	var o FormatOsu
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return &o, err
	}
	dat = bytes.ReplaceAll(dat, []byte("\r\n"), []byte("\n"))

	o.Events = make([]Event, 0)
	o.TimingPoints = make([]TimingPoint, 0)
	o.HitObjects = make([]HitObject, 0)
	var section string
	for _, l := range bytes.Split(dat, []byte("\n")) {
		l = bytes.TrimSpace(l)
		line := string(l)
		if isPass(line) {
			continue
		}
		if isSection(line) {
			section = strings.Trim(line, "[]")
			continue
		}
		switch section {
		// %s
		case "Events":
			continue
		case "TimingPoints":
			tp, err := newTimingPoint(line)
			if err != nil {
				return &o, err
			}
			o.TimingPoints = append(o.TimingPoints, tp)
		}
	}
	return &o, nil
}
func isPass(line string) bool {
	return len(line) == 0 || len(line) >= 2 && line[:2] == "//"
}
func isSection(line string) bool {
	if len(line) == 0 {
		return false
	}
	return string(line[0]) == "[" && string(line[len(line)-1]) == "]"
}
