package piano

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/mode"
)

// Chart should avoid redundant data as much as possible
type Chart struct {
	mode.Chart
	// mode.ChartHeader
	// TransPoints []*mode.TransPoint
	// Mode        int
	KeyCount int
	Notes    []Note
	// Duration    int64
	// NoteCounts  []int
}

// 7 key chart's Mode is 128 + 7 = 135
func NewChart(fpath string) (*Chart, error) { // f any
	var c Chart
	dat, err := os.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	var f any
	switch strings.ToLower(filepath.Ext(fpath)) {
	case ".osu":
		f, err = osu.Parse(dat)
		if err != nil {
			return nil, err
		}
	}
	c.ChartHeader = mode.NewChartHeader(f)
	c.TransPoints = mode.NewTransPoints(f)

	switch f := f.(type) {
	case *osu.Format:
		c.KeyCount = int(f.CircleSize)
		if c.KeyCount <= 4 {
			c.Mode = mode.ModePiano4
		} else {
			c.Mode = mode.ModePiano7
		}
		// c.Mode = mode.ModePiano
		// c.Mode += c.KeyCount
		c.Mode2 = c.KeyCount
		c.Notes = make([]Note, 0, len(f.HitObjects)*2)
		for _, ho := range f.HitObjects {
			c.Notes = append(c.Notes, NewNote(ho, c.KeyCount)...)
		}
	}
	sort.Slice(c.Notes, func(i, j int) bool {
		if c.Notes[i].Time == c.Notes[j].Time {
			return c.Notes[i].Key < c.Notes[j].Key
		}
		return c.Notes[i].Time < c.Notes[j].Time
	})
	if len(c.Notes) > 0 {
		c.Duration = c.Notes[len(c.Notes)-1].Time
	}
	c.NoteCounts = make([]int, 2)
	for _, n := range c.Notes {
		switch n.Type {
		case Normal:
			c.NoteCounts[0]++
		case Head:
			c.NoteCounts[1]++
		}
	}
	return &c, nil
}
