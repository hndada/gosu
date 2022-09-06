package drum

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu"
	"github.com/hndada/gosu/format/osu"
)

type Floater struct {
	Time  int64
	Speed float64
}

func (f Floater) Position(time int64) float64 {
	return float64(f.Time-time) * f.Speed
}

type Chart struct {
	gosu.ChartHeader
	TransPoints []*gosu.TransPoint
	Notes       []*Note
	Dots        []*Dot
	Bars        []*Bar
}

var (
	Showtime    int64   = 500
	TickDensity float64 = 4.0
	// ShakeDensity = 4.0 // One beat has 4 Shakes.
)

// NewChart takes file path as input for starting with parsing.
// Chart data should not rely on the ChartInfo; users may have modified it.
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
	mainBPM, _, _ := c.BPMs()
	bpmScale := c.TransPoints[0].BPM / mainBPM
	for _, tp := range c.TransPoints {
		tp.Speed *= bpmScale
	}

	c.Notes = NewNotes(f)
	tp := c.TransPoints[0]
	for _, n := range c.Notes {
		for tp.Next != nil && n.Time >= tp.Next.Time {
			tp = tp.Next
		}
		n.Speed = tp.Speed
		bpm := ScaledBPM(tp.BPM)
		switch f := f.(type) {
		case *osu.Format:
			// speedFactor := c.TransPoints[0].BPM / 60000 * (f.SliderMultiplier * 100)
			speed := tp.BPM * (tp.Speed / bpmScale) * f.SliderMultiplier * 100
			n.Duration = int64(n.length / speed)
		}
		n.Tick = int(float64(n.Duration) / bpm * TickDensity)
	}
	c.Dots = NewDots(c.Notes)
	// switch f := f.(type) {
	// case *osu.Format:
	// 	// TransPoints' speed has not scaled yet.
	// 	tp := c.TransPoints[0]
	// 	for _, n := range c.Notes {
	// 		for tp.Next != nil && n.Time >= tp.Next.Time {
	// 			tp = tp.Next
	// 		}

	// 	}
	// }
	c.Bars = NewBars(c.TransPoints, c.Duration())
	tp = c.TransPoints[0]
	for _, b := range c.Bars {
		for tp.Next != nil && b.Time >= tp.Next.Time {
			tp = tp.Next
		}
		b.Speed = tp.Speed
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
	last := c.Notes[len(c.Notes)-1]
	return last.Time + last.Duration
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

// It is proved that all BPMs are set into [MinScaledBPM, MaxScaledBPM) by v*2 or v/2
// if MinScaledBPM *2 >= MaxScaleBPM.
func ScaledBPM(bpm float64) float64 {
	if bpm < 0 {
		bpm = -bpm
	}
	switch {
	case bpm > MaxScaledBPM:
		for bpm > MaxScaledBPM {
			bpm /= 2
		}
	case bpm >= MinScaledBPM:
		return bpm
	case bpm < MinScaledBPM:
		for bpm < MinScaledBPM {
			bpm *= 2
		}
	}
	return bpm
}
