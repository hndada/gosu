package drum

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/mode"
)

const Mode = 1 // ModeDrum
type Floater struct {
	Time  int64
	Speed float64
}

func (f Floater) Position(time int64) float64 {
	return float64(f.Time-time) * f.Speed
}

type Chart struct {
	mode.ChartHeader
	MD5         [16]byte
	TransPoints []*mode.TransPoint
	Notes       []*Note
	Rolls       []*Note
	Shakes      []*Note
	Dots        []*Dot // Ticks in a Roll note.
	Bars        []*Bar

	Level        float64
	ScoreFactors [3]float64
}

var (
	// TickDensity  float64 = 4
	DotDensity   float64 = 4 // Infers how many dots per beat in Roll note.
	ShakeDensity float64 = 3 // Infers how many shakes per beat in Shake note.
)

// NewChart takes file path as input for starting with parsing.
// Chart data should not rely on the ChartInfo; users may have modified it.
func NewChart(fsys fs.FS, name string) (c *Chart, err error) {
	var dat []byte
	dat, err = fs.ReadFile(fsys, name)
	if err != nil {
		name = strings.Replace(name, "[4K] ", "", 1)
		name = strings.Replace(name, "[7K] ", "", 1)
		err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if strings.Contains(d.Name(), name) {
				// fmt.Println("found", d.Name(), name)
				dat, err = fs.ReadFile(fsys, d.Name())
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	name += ".osu"
	var f any
	switch filepath.Ext(name) {
	case ".osu":
		f, err = osu.Parse(dat)
		if err != nil {
			return
		}
	}
	c = new(Chart)
	c.ChartHeader = mode.NewChartHeader(f)
	c.Mode = Mode
	c.SubMode = 4
	c.MD5 = md5.Sum(dat)
	c.TransPoints = mode.NewTransPoints(f)
	if len(c.TransPoints) == 0 {
		err = fmt.Errorf("no TransPoints in the chart")
		return
	}
	mainBPM, _, _ := c.BPMs()
	bpmScale := c.TransPoints[0].BPM / mainBPM
	for _, tp := range c.TransPoints {
		tp.Speed *= bpmScale
	}

	c.Notes, c.Rolls, c.Shakes = NewNotes(f)
	var tp *mode.TransPoint
	for _, ns := range [][]*Note{c.Notes, c.Rolls, c.Shakes} {
		tp = c.TransPoints[0]
		for _, n := range ns {
			for tp.Next != nil && n.Time >= tp.Next.Time {
				tp = tp.Next
			}
			n.Speed = tp.Speed
			bpm := ScaledBPM(tp.BPM)
			if n.Type == Roll {
				switch f := f.(type) {
				case *osu.Format:
					// speedFactor := c.TransPoints[0].BPM / 60000 * (f.SliderMultiplier * 100)
					speed := tp.BPM * (tp.Speed / bpmScale) / 60000 * f.SliderMultiplier * 100 // Unit is osupixel / 100ms.
					n.Duration = int64(n.length / speed)
				}
			}
			switch n.Type {
			case Roll:
				n.Tick = int(float64(n.Duration)*bpm/60000*DotDensity+0.1) + 1
			case Shake:
				n.Tick = int(float64(n.Duration)*bpm/60000*ShakeDensity+0.1) + 1
			}
		}
	}
	c.Dots = NewDots(c.Rolls)
	c.Bars = NewBars(c.TransPoints, c.Duration())
	tp = c.TransPoints[0]
	for _, b := range c.Bars {
		for tp.Next != nil && b.Time >= tp.Next.Time {
			tp = tp.Next
		}
		b.Speed = tp.Speed
	}
	c.Level, c.ScoreFactors = mode.Level(c)
	return
}

const (
	MaxScaledBPM = 280 // 256
	MinScaledBPM = 60  // 128
)

func (c Chart) Duration() (last int64) {
	for _, ns := range [][]*Note{c.Notes, c.Rolls, c.Shakes} {
		if len(ns) == 0 {
			continue
		}
		n := ns[len(ns)-1]
		if last2 := n.Time + n.Duration; last < last2 {
			last = last2
		}
	}
	return
}
func (c Chart) NoteCounts() (vs []int) {
	vs = make([]int, 3)
	for _, n := range c.Notes {
		vs[n.Type]++
	}
	return
}
func (c Chart) BPMs() (main, min, max float64) {
	return mode.BPMs(c.TransPoints, c.Duration())
}

// It is proved that all BPMs are set into [min, max) by v*2 or v/2 if 2 * min >= max.
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

// func NewChartInfo(cpath string) (info game.ChartInfo, err error) {
// 	c, err := NewChart(cpath)
// 	if err != nil {
// 		return
// 	}
// 	// mode := game.ModeDrum
// 	main, min, max := c.BPMs()
// 	info = game.ChartInfo{
// 		Path: cpath,
// 		// Mods:       mods,
// 		ChartHeader: c.ChartHeader,
// 		// Mode:        mode,
// 		// SubMode:     0,
// 		Level:      c.Level,
// 		Duration:   c.Duration(),
// 		NoteCounts: c.NoteCounts(),
// 		MainBPM:    main,
// 		MinBPM:     min,
// 		MaxBPM:     max,
// 	}
// 	return
// }
