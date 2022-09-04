package piano

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/format/osu"
)

type Chart struct {
	gosu.ChartHeader
	KeyCount    int
	TransPoints []*gosu.TransPoint
	Notes       []*Note
	Bars        []*Bar
}

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
	switch f := f.(type) {
	case *osu.Format:
		c.KeyCount = int(f.CircleSize)
	}
	c.Notes = NewNotes(f, c.KeyCount)
	c.Bars = NewBars(c.TransPoints, c.Duration())

	// Position calculation is based on TransPoints.
	mainBPM, _, _ := c.BPMs()
	bpmScale := c.TransPoints[0].BPM / mainBPM
	for _, tp := range c.TransPoints {
		tp.Speed *= bpmScale
		if prev := tp.Prev; prev != nil {
			tp.Position = prev.Position + float64(tp.Time-prev.Time)*prev.Speed
		} else {
			tp.Position = float64(tp.Time) * tp.Speed
		}
	}
	tp := c.TransPoints[0]
	for _, n := range c.Notes {
		for tp.Next != nil && n.Time >= tp.Next.Time {
			tp = tp.Next
		}
		n.Position = tp.Position + float64(n.Time-tp.Time)*tp.Speed
	}
	tp = c.TransPoints[0]
	for _, b := range c.Bars {
		for tp.Next != nil && b.Time >= tp.Next.Time {
			tp = tp.Next
		}
		b.Position = tp.Position + float64(b.Time-tp.Time)*tp.Speed
	}
	return
}

func (c Chart) Duration() int64 {
	if len(c.Notes) == 0 {
		return 0
	}
	return c.Notes[len(c.Notes)-1].Time
}
func (c Chart) NoteCounts() (vs []int) {
	vs = make([]int, 2)
	for _, n := range c.Notes {
		switch n.Type {
		case Normal:
			vs[0]++
		case Head:
			vs[1]++
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
	mode := gosu.ModePiano4
	if c.KeyCount > 4 {
		mode = gosu.ModePiano7
	}
	main, min, max := c.BPMs()
	info = gosu.ChartInfo{
		Path: cpath,
		// Mods:       mods,
		Header:     c.ChartHeader,
		Mode:       mode,
		SubMode:    c.KeyCount,
		Level:      gosu.Level(c),
		Duration:   c.Duration(),
		NoteCounts: c.NoteCounts(),
		MainBPM:    main,
		MinBPM:     min,
		MaxBPM:     max,
	}
	return
}
