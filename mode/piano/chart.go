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

	Mods     interface{}
	Dynamics []*mode.Dynamic
	Notes    []*Note
	Bars     []*Bar
}

func LoadChart(fsys fs.FS, name string, mods interface{}) (c *Chart, err error) {
	format, hash, err := mode.ParseChartFile(fsys, name)
	if err != nil {
		return
	}

	c = new(Chart)
	c.ChartHeader = mode.LoadChartHeader(format)
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
	c.setNotePositions()
	c.setBarPositions()
	return
}

// Position is for drawing notes and bars efficiently.
// Only cursor is updated in every Update(), then notes and bars
// are drawn based on the difference between their positions and cursor's.
// Position calculation is based on Dynamics.
func (c *Chart) setDynamicPositions() {
	mainBPM, _, _ := mode.BPMs(c.Dynamics, c.Duration())
	bpmScale := c.Dynamics[0].BPM / mainBPM
	for _, dy := range c.Dynamics {
		dy.Speed *= bpmScale
		if prev := dy.Prev; prev != nil {
			dy.Position = prev.Position + float64(dy.Time-prev.Time)*prev.Speed
		} else {
			dy.Position = float64(dy.Time) * dy.Speed
		}
	}
}
func (c *Chart) setNotePositions() {
	dy := c.Dynamics[0]
	for _, n := range c.Notes {
		for dy.Next != nil && n.Time >= dy.Next.Time {
			dy = dy.Next
		}
		n.Position = dy.Position + float64(n.Time-dy.Time)*dy.Speed
		if n.Type == Tail {
			n.Position += float64(S.TailExtraTime) * dy.Speed

			// Tail notes should be drawn after their heads.
			if n.Position < n.Prev.Position {
				n.Position = n.Prev.Position
			}
		}
	}
}
func (c *Chart) setBarPositions() {
	dy := c.Dynamics[0]
	for _, b := range c.Bars {
		for dy.Next != nil && b.Time >= dy.Next.Time {
			dy = dy.Next
		}
		b.Position = dy.Position + float64(b.Time-dy.Time)*dy.Speed
	}
}

func (c Chart) Duration() int64 {
	if len(c.Notes) == 0 {
		return 0
	}
	last := c.Notes[len(c.Notes)-1]
	return last.Time + last.Duration
}
