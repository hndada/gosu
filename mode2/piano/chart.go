package piano

import (
	"io/fs"

	mode "github.com/hndada/gosu/mode2"
)

type Chart struct {
	mode.ChartHeader
	Notes []*Note
	mode.Dynamics
	Bars []*mode.Bar
	Level
}

// NewXxx returns *Chart, while LoadXxx doesn't.
func NewChart(fsys fs.FS, name string) (c *Chart, err error) {
	format, hash, err := mode.LoadChartFile(fsys, name)
	if err != nil {
		return
	}

	c.ChartHeader = mode.NewChartHeader(format, hash)
	c.Notes = NewNotes(format, c.KeyCount())
	c.Dynamics, err = mode.NewDynamics(format, c.Duration())
	if err != nil {
		return
	}
	c.Bars = c.Dynamics.NewBars(c.Duration())

	c.setDynamicPositions()
	c.setNotePositions()
	c.setBarPositions()

	c.setSteps()
	c.setLevel()
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
	mainBPM, _, _ := c.Dynamics.BPMs(c.Duration())
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

func (c *Chart) setNotePositions() {
	d := c.Dynamics[0]
	for _, n := range c.Notes {
		for d.Next != nil && n.Time >= d.Next.Time {
			d = d.Next
		}
		n.Position = d.Position + float64(n.Time-d.Time)*d.Speed

		if n.Type == Tail {
			// Tail notes should be drawn after their heads.
			if n.Position < n.Prev.Position {
				n.Position = n.Prev.Position
			}
		}
	}
}

func (c *Chart) updateTailPosition(tailExtraDuration int32) {
	d := c.Dynamics[0]
	for _, n := range c.Notes {
		if n.Type != Tail {
			continue
		}
		for d.Next != nil && n.Time >= d.Next.Time {
			d = d.Next
		}

		n.Position += float64(tailExtraDuration) * d.Speed
		// Tail notes should be drawn after their heads.
		if n.Position < n.Prev.Position {
			n.Position = n.Prev.Position
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

func (c Chart) newStagedNotes() []*Note {
	staged := make([]*Note, c.KeyCount())
	for k := range staged {
		for _, n := range c.Notes {
			if k == n.Key {
				staged[n.Key] = n
				break
			}
		}
	}
	return staged
}

func (c Chart) NoteCounts() []int {
	counts := make([]int, 3)
	for _, n := range c.Notes {
		switch n.Type {
		case Normal:
			counts[0]++
		case Head:
			counts[1]++
		case Tail:
			counts[2]++
		}
	}
	return counts
}

func (c Chart) KeyCount() int { return c.SubMode }
