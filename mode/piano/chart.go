package piano

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/mode"
)

type Chart struct {
	mode.ChartHeader
	KeyCount int      // same with ChartHeader.SubMode
	Hash     [16]byte // MD5

	Mods     Mods
	Dynamics []*mode.Dynamic
	Notes    []*Note
	Bars     []*Bar
}

// NewXxx returns *Chart, while LoadXxx doesn't.
func NewChart(cfg *Config, fsys fs.FS, name string, mods Mods) (c *Chart, err error) {
	c = new(Chart)

	format, hash, err := mode.ParseChartFile(fsys, name)
	if err != nil {
		return
	}
	c.ChartHeader = mode.NewChartHeader(format)
	c.KeyCount = c.SubMode
	c.Hash = hash

	c.Mods = mods
	c.Dynamics = mode.NewDynamics(format)
	if len(c.Dynamics) == 0 {
		err = fmt.Errorf("no Dynamics in the chart")
		return
	}
	c.Notes = NewNotes(format, c.KeyCount)
	c.Bars = NewBars(c.Dynamics, c.Duration())

	c.setDynamicPositions()
	c.setNotePositions(cfg)
	c.setBarPositions()
	return
}

// Position is for drawing notes and bars efficiently.
// Only cursor is updated in every Update(), then notes and bars
// are drawn based on the difference between their positions and cursor's.

// Position calculation is based on Dynamics.
// Farther note has larger position.
// Tail's Position is always larger than Head's.
func (c *Chart) setDynamicPositions() {
	// Brilliant idea: Make SpeedScale scaled by MainBPM.
	mainBPM, _, _ := mode.BPMs(c.Dynamics, c.Duration())
	bpmScale := c.Dynamics[0].BPM / mainBPM
	for _, d := range c.Dynamics {
		d.Speed *= bpmScale
		if prev := d.Prev; prev != nil {
			d.Position = prev.Position + float64(d.Time-prev.Time)*prev.Speed
		} else {
			d.Position = float64(d.Time) * d.Speed
		}
	}
}
func (c *Chart) setNotePositions(cfg *Config) {
	tailExtraTime := cfg.TailExtraDuration
	d := c.Dynamics[0]
	for _, n := range c.Notes {
		for d.Next != nil && n.Time >= d.Next.Time {
			d = d.Next
		}
		n.Position = d.Position + float64(n.Time-d.Time)*d.Speed
		if n.Type == Tail {
			n.Position += float64(tailExtraTime) * d.Speed

			// Tail notes should be drawn after their heads.
			if n.Position < n.Prev.Position {
				n.Position = n.Prev.Position
			}
		}
	}
}
func (c *Chart) setBarPositions() {
	d := c.Dynamics[0]
	for _, b := range c.Bars {
		for d.Next != nil && b.Time >= d.Next.Time {
			d = d.Next
		}
		b.Position = d.Position + float64(b.Time-d.Time)*d.Speed
	}
}

func (c Chart) Duration() int32 {
	if len(c.Notes) == 0 {
		return 0
	}
	last := c.Notes[len(c.Notes)-1]
	return last.Time // + last.Duration
}
