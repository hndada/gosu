package piano

import (
	"io/fs"

	"github.com/hndada/gosu/game"
)

// Todo: make fields unexported?
type Chart struct {
	Mods Mods
	game.ChartHeader
	game.Dynamics
	Notes
	// KeyCount int
}

func NewChart(fsys fs.FS, name string, mods Mods) (*Chart, error) {
	c := &Chart{
		Mods: mods,
	}

	format, hash, err := game.LoadChartFormat(fsys, name)
	if err != nil {
		return c, err
	}
	c.ChartHeader = game.NewChartHeader(format, hash)
	// c.KeyCount = c.SubMode

	dys, err := game.NewDynamics(format)
	if err != nil {
		return c, err
	}
	c.Dynamics = dys

	keyCount := c.SubMode
	c.Notes = NewNotes(keyCount, format, dys)
	return c, nil
}

func (c Chart) NoteCounts() []int {
	counts := make([]int, 2)
	for _, n := range c.Notes.data {
		switch n.Kind {
		case Normal:
			counts[0]++
		case Head:
			counts[1]++
		}
	}
	return counts
}

func (c Chart) TotalDuration() int32 {
	ns := c.Notes.data
	if len(ns) == 0 {
		return 0
	}

	// No need to add last.Duration, since last is
	// always either Normal or Tail.
	last := ns[len(ns)-1]
	return last.Time
}
