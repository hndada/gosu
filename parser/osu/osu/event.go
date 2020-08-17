package osu

import (
	"strconv"
	"strings"
)

func newEvent(line string) (Event, error) {
	vs := strings.Split(line, `,`)
	var e Event
	switch vs[0] {
	case "0", "1", "Video":
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
	case "2", "Break":
		{
			f, err := strconv.ParseFloat(vs[1], 64)
			if err != nil {
				return e, err
			}
			e.StartTime = int(f)
		}
		{
			f, err := strconv.ParseFloat(vs[1], 64)
			if err != nil {
				return e, err
			}
			e.EndTime = int(f)
		}
	}
	return e, nil
}
