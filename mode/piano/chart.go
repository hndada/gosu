package piano

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/mode"
)

const Mode = 0 // ModePiano
// Level, ScoreFactors, MD5 will not exported to file.
type Chart struct {
	mode.ChartHeader
	MD5  [16]byte
	Mods interface{}

	KeyCount    int // The number of physical keys.
	KeyMode     int // Scratch mode value may add after applying mods.
	TransPoints []*mode.TransPoint
	Notes       []*Note
	Bars        []*Bar

	Level        float64
	ScoreFactors [3]float64
}

// Position is for calculating note and bar's sprite positions efficiently.
// Positions of notes and bars at time = 0 are calculated in advance.
// In every Update(), only current cursor's Position is calculated.
// Notes and bars are drawn based on the difference between their positions and cursor's.
func NewChart(fsys fs.FS, name string) (c *Chart, err error) {
	var dat []byte
	dat, err = fs.ReadFile(fsys, name)
	if err != nil {
		return
	}
	var f any // f stands for Format. // Todo: f -> _type?
	switch filepath.Ext(name) {
	case ".osu", ".OSU":
		f, err = osu.Parse(dat)
		if err != nil {
			return
		}
	}
	c = new(Chart)
	c.ChartHeader = mode.NewChartHeader(f)
	c.Mode = Mode
	switch f := f.(type) {
	case *osu.Format:
		c.SubMode = int(f.CircleSize)
	}
	c.MD5 = md5.Sum(dat)
	switch f := f.(type) {
	case *osu.Format:
		c.KeyCount = int(f.CircleSize)
		c.KeyMode = c.KeyCount
	}
	c.TransPoints = mode.NewTransPoints(f)
	if len(c.TransPoints) == 0 {
		err = fmt.Errorf("no TransPoints in the chart")
		return
	}
	c.Notes = NewNotes(f, c.KeyMode)
	c.Bars = NewBars(c.TransPoints, c.Duration())

	// Calculate positions. Position calculation is based on TransPoints.
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
		if n.Type == Tail {
			n.Position += float64(S.TailExtraTime) * tp.Speed
			// It is guaranteed that Tail's Prev (which is Head)
			// has already proceeded, since c.Notes is sorted by Time.
			if n.Position < n.Prev.Position {
				n.Position = n.Prev.Position
			}
		}
	}
	tp = c.TransPoints[0]
	for _, b := range c.Bars {
		for tp.Next != nil && b.Time >= tp.Next.Time {
			tp = tp.Next
		}
		b.Position = tp.Position + float64(b.Time-tp.Time)*tp.Speed
	}
	c.Level, c.ScoreFactors = mode.Level(c)
	return
}

func (c Chart) Duration() int64 {
	if len(c.Notes) == 0 {
		return 0
	}
	last := c.Notes[len(c.Notes)-1]
	return last.Time + last.Duration
}

//	func (c Chart) NoteCounts() (vs []int) {
//		vs = make([]int, 2)
//		for _, n := range c.Notes {
//			if n.Type == Tail {
//				continue
//			}
//			vs[n.Type]++
//		}
//		return
//	}
func (c Chart) BPMs() (main, min, max float64) {
	return mode.BPMs(c.TransPoints, c.Duration())
}

// func (c Chart) NoteCountString() string {
// 	vs := c.NoteCounts()
// 	total := vs[0] + 2*vs[1]
// 	ratio := float64(vs[0]) / float64(vs[0]+vs[1])
// 	return fmt.Sprintf("Notes: %d\nLN: %.0f%%", total, ratio*100)
// }

// func NewChartInfo(cpath string) (info game.ChartInfo, err error) {
// 	c, err := NewChart(cpath)
// 	if err != nil {
// 		return
// 	}
// 	// mode := game.ModePiano4
// 	// if c.KeyCount > 4 {
// 	// 	mode = game.ModePiano7
// 	// }
// 	main, min, max := c.BPMs()
// 	info = game.ChartInfo{
// 		Path: cpath,
// 		// Mods:       mods,
// 		ChartHeader: c.ChartHeader,
// 		// Mode:        mode,
// 		// SubMode:     c.KeyCount,
// 		Level:      c.Level,
// 		Duration:   c.Duration(),
// 		NoteCounts: c.NoteCounts(),
// 		MainBPM:    main,
// 		MinBPM:     min,
// 		MaxBPM:     max,
// 	}
// 	return
// }
