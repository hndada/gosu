package drum

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/format/osu"
)

type Chart struct {
	gosu.BaseChart
	Notes []Note
}

// NewChart takes file path as input for starting with parsing.
// Chart data should not rely on the ChartInfo; clients may have modified it.
func NewChart(cpath string, mods gosu.Mods) (*Chart, error) {
	var c Chart
	dat, err := os.ReadFile(cpath)
	if err != nil {
		return nil, err
	}
	var f any
	switch strings.ToLower(filepath.Ext(cpath)) {
	case ".osu":
		f, err = osu.Parse(dat)
		if err != nil {
			return nil, err
		}
		// osu.Parse doesn't sort Hit objects, not to manipulate unprompted tasks.
		// Following tasks (adding notes to Chart) assume timing points and
		// hit objects are sorted well, otherwise unexpected result would return.
		// Glad to know that gosu.NewTransPoints() does sort in the task.

		// The Verifier (check whether next hit object's time is equal or later
		// than previous one) might be possible to be added, but it still requires
		// some computation time (which is O(n)).
	}
	c.ChartHeader = gosu.NewChartHeader(f)
	c.TransPoints = gosu.NewTransPoints(f)
	c.ModeType = gosu.ModeTypeDrum
	// No sub mode for Drum mode.
	switch f := f.(type) {
	case *osu.Format:
		c.Notes = make([]Note, 0, len(f.HitObjects)*2)
		if len(c.TransPoints) == 0 {
			return nil, fmt.Errorf("no TransPoints in the chart")
		}
		tp := c.TransPoints[0]
		for _, ho := range f.HitObjects {
			for tp.Next != nil && tp.Next.Time <= int64(ho.Time) {
				tp = tp.Next
			}
			speed := tp.BPM * tp.BeatLengthScale * f.SliderMultiplier
			scaledBPM := ScaledBPM(tp.BPM)
			c.Notes = append(c.Notes, NewNote(ho, speed, scaledBPM)...)
		}
	}
	sort.Slice(c.Notes, func(i, j int) bool {
		if c.Notes[i].Time == c.Notes[j].Time {
			return c.Notes[i].Type < c.Notes[j].Type
		}
		return c.Notes[i].Time < c.Notes[j].Time
	})
	if len(c.Notes) > 0 {
		c.Duration = c.Notes[len(c.Notes)-1].Time
	}
	c.NoteCounts = make([]int, 3)
	for _, n := range c.Notes {
		switch n.Type {
		case Don, Kat, BigDon, BigKat:
			c.NoteCounts[0]++
		case Head, BigHead:
			c.NoteCounts[1]++
		case Shake:
			c.NoteCounts[2]++
		}
	}
	return &c, nil
}

func NewChartInfo(cpath string, mods gosu.Mods) (gosu.ChartInfo, error) {
	c, err := NewChart(cpath, mods)
	if err != nil {
		return gosu.ChartInfo{}, err
	}
	return gosu.NewChartInfo(&c.BaseChart, cpath, gosu.Level(c)), nil
}
