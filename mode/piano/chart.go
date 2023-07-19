package piano

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/mode"
)

type Chart struct {
	mode.ChartHeader
	KeyCount int // same with ChartHeader.SubMode

	Dynamics []*mode.Dynamic
	Notes    []*Note
	Bars     []*Bar
}

// NewXxx returns *Chart, while LoadXxx doesn't.
func NewChart(cfg *Config, fsys fs.FS, name string) (*Chart, error) {
	c := new(Chart)
	f, err := fsys.Open(name)
	if err != nil {
		return c, fmt.Errorf("open %s: %w", name, err)
	}

	var format any
	switch filepath.Ext(name) {
	case ".osu", ".OSU":
		format, err = osu.NewFormat(f)
		if err != nil {
			return c, fmt.Errorf("new osu format: %w", err)
		}
	}

	c.ChartHeader = mode.NewChartHeader(format)
	c.KeyCount = c.SubMode
	c.ChartHash, _ = mode.Hash(f)

	c.Dynamics = mode.NewDynamics(format)
	if len(c.Dynamics) == 0 {
		return c, fmt.Errorf("no Dynamics in the chart")
	}
	c.Notes = NewNotes(format, c.KeyCount)
	c.Bars = NewBars(c.Dynamics, c.Duration())

	c.setDynamicPositions()
	c.setNotePositions(cfg)
	c.setBarPositions()
	return c, nil
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
