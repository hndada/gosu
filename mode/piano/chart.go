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
	TransPoints  []*gosu.TransPoint
	KeyCount     int
	SpeedScale   float64 // Affects Note and Bar's position.
	Notes        []*Note
	BarPositions []float64
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
	c.ChartHeader = gosu.NewChartHeader(f)
	fixed := true
	c.TransPoints = gosu.NewTransPoints(f, fixed)
	switch f := f.(type) {
	case *osu.Format:
		c.KeyCount = int(f.CircleSize)
	}
	c.SpeedScale = 1
	c.Notes = NewNotes(f, c.KeyCount, c.TransPoints[0])
	c.BarPositions = make([]float64, 0)
	return
}
func (c *Chart) SetSpeedScale(speedScale float64) {
	for _, n := range c.Notes {
		n.Position *= speedScale / c.SpeedScale
	}
	c.SpeedScale = speedScale
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
	main, min, max := gosu.BPMs(c.TransPoints, c.Duration())
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
