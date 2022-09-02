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
	// SpeedScale  float64 // Affects Note and Bar's position.
	// MainBPM     float64
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
	fixed := true
	c.TransPoints = gosu.NewTransPoints(f, fixed)
	switch f := f.(type) {
	case *osu.Format:
		c.KeyCount = int(f.CircleSize)
	}
	c.Notes = NewNotes(f, c.KeyCount, c.TransPoints[0])
	c.Bars = NewBars(c.TransPoints, c.Duration())
	// c.SpeedScale = 1
	// c.MainBPM, _, _ = c.BPMs()
	// mainBPM, _, _ := c.BPMs()
	// for _, tp := range c.TransPoints {
	// 	tp.Position /= mainBPM
	// }
	// for _, n := range c.Notes {
	// 	n.Position /= mainBPM
	// }
	// for _, b := range c.Bars {
	// 	b.Position /= mainBPM
	// }
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
