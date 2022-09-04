package drum

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/format/osu"
)

type Chart struct {
	gosu.ChartHeader
	TransPoints []*gosu.TransPoint
	Notes       []*Note
	Bars        []*Bar
}

// NewChart takes file path as input for starting with parsing.
// Chart data should not rely on the ChartInfo; clients may have modified it.

// Position is for calculating note and bar's sprite positions efficiently.
// Positions of notes and bars at time = 0 are calculated in advance.
// In every Update(), only current cursor's Position is calculated.
// Notes and bars are drawn based on the difference between their positions and cursor's.
func NewChart(cpath string) (c *Chart, err error) {
	var f any
	dat, err := os.ReadFile(cpath)
	if err != nil {
		return
	}
	switch strings.ToLower(filepath.Ext(cpath)) {
	case ".osu":
		f, err = osu.Parse(dat)
		if err != nil {
			return
		}
	}
	c = new(Chart)
	c.ChartHeader = gosu.NewChartHeader(f)
	c.TransPoints = gosu.NewTransPoints(f)
	if len(c.TransPoints) == 0 {
		err = fmt.Errorf("no TransPoints in the chart")
		return
	}
	c.Notes = NewNotes(f, c.TransPoints)
	c.Bars = NewBars(c.TransPoints, c.Duration())

	// Calculate positions and speed.
	// Position calculation is based on TransPoints.
	mainBPM, _, _ := c.BPMs()
	bpmScale := c.TransPoints[0].BPM / mainBPM
	for _, tp := range c.TransPoints {
		tp.Speed *= bpmScale
		tp.Position = float64(tp.Time) * tp.Speed
	}
	tp := c.TransPoints[0]
	for _, n := range c.Notes {
		for tp.Next != nil && n.Time >= tp.Next.Time {
			tp = tp.Next
		}
		n.Speed = tp.Speed
		n.Position = float64(n.Time) * n.Speed
	}
	tp = c.TransPoints[0]
	for _, b := range c.Bars {
		for tp.Next != nil && b.Time >= tp.Next.Time {
			tp = tp.Next
		}
		b.Speed = tp.Speed
		b.Position = float64(b.Time) * b.Speed
	}
	return
}

const (
	MaxScaledBPM = 280 // 256
	MinScaledBPM = 60  // 128
)

func (c Chart) Duration() int64 {
	if len(c.Notes) == 0 {
		return 0
	}
	return c.Notes[len(c.Notes)-1].Time
}
func (c Chart) NoteCounts() (vs []int) {
	vs = make([]int, 3)
	for _, n := range c.Notes {
		switch n.Type {
		case Normal:
			vs[0]++
		case Head:
			vs[1]++
		case Shake:
			vs[2]++
		}
	}
	return
}
func (c Chart) BPMs() (main, min, max float64) {
	return gosu.BPMs(c.TransPoints, c.Duration())
}
func NewChartInfo(cpath string) (info gosu.ChartInfo, err error) {
	c, err := NewChart(cpath)
	if err != nil {
		return
	}
	// Todo: put mods implementation here
	mode := gosu.ModeDrum
	main, min, max := c.BPMs()
	info = gosu.ChartInfo{
		Path: cpath,
		// Mods:       mods,
		Header:     c.ChartHeader,
		Mode:       mode,
		SubMode:    0,
		Level:      gosu.Level(c),
		Duration:   c.Duration(),
		NoteCounts: c.NoteCounts(),
		MainBPM:    main,
		MinBPM:     min,
		MaxBPM:     max,
	}
	return
}
