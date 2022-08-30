package piano

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/format/osu"
)

// Chart should avoid redundant data as much as possible
type Chart struct {
	gosu.BaseChart
	KeyCount int
	Notes    []*gosu.Note
}

// NewChart takes file path as input for starting with parsing.
// Chart data should not rely on the ChartInfo; clients may have compromised it.
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
	}

	c.ChartHeader = gosu.NewChartHeader(f)
	c.TransPoints = gosu.NewTransPoints(f)
	mainBPM, _, _ := gosu.BPMs(c.TransPoints, c.Duration)
	// In Piano mode, TransPoint's Position should be divided by main BPM for scaling.
	for i := range c.TransPoints {
		c.TransPoints[i].Position /= mainBPM
	}
	switch f := f.(type) {
	case *osu.Format:
		c.KeyCount = int(f.CircleSize)
		if c.KeyCount <= 4 {
			c.ModeType = gosu.ModeTypePiano4
		} else {
			c.ModeType = gosu.ModeTypePiano7
		}
		c.SubModeType = c.KeyCount
		c.Notes = make([]*gosu.Note, 0, len(f.HitObjects)*2)
		for _, ho := range f.HitObjects {
			// bns := gosu.NewNote(ho)
			// for _, bn := range bns {
			// 	c.Notes = append(c.Notes, Note{
			// 		Note: bn,
			// 		Key:  ho.Column(c.SubMode),
			// 	})
			// }
			c.Notes = append(c.Notes, gosu.NewNote(ho, c.ModeType, c.SubModeType)...)
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
	c.SetBars()
	c.NoteCounts = make([]int, 2)
	for _, n := range c.Notes {
		switch n.Type {
		case gosu.Normal:
			c.NoteCounts[0]++
		case gosu.Head:
			c.NoteCounts[1]++
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
