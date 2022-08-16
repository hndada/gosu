package piano

import (
	"sort"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/mode"
)

// Chart should avoid redundant data as much as possible
type Chart struct {
	mode.Chart
	KeyCount int
	Notes    []Note
	// ScratchMode int
}

// Todo: remove return value error
func NewChartFromOsu(o *osu.Format) (*Chart, error) {
	c := &Chart{
		Chart:    mode.NewChartFromOsu(o),
		KeyCount: int(o.CircleSize),
		Notes:    make([]Note, 0, len(o.HitObjects)*2),
	}
	for _, ho := range o.HitObjects {
		c.Notes = append(c.Notes, NewNoteFromOsu(ho, c.KeyCount)...)
	}
	sort.Slice(c.Notes, func(i, j int) bool {
		if c.Notes[i].Time == c.Notes[j].Time {
			return c.Notes[i].Key < c.Notes[j].Key
		}
		return c.Notes[i].Time < c.Notes[j].Time
	})
	c.EndTime = func() int64 { return c.Notes[len(c.Notes)-1].Time }
	return c, nil
}
