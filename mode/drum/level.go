package drum

import (
	"math"
)

func (n Note) Weight() float64 {
	switch n.Type {
	case Normal:
		switch n.Size {
		case Regular:
			return 1.0
		case Big:
			return 1.1
		}
	case Shake:
		// Shake is apparently easier than Roll, since it is free from beat.
		// https://www.desmos.com/calculator/nsogcrebx9
		return math.Pow(float64(32*n.Tick), 0.75) / 32 // 32 comes from 8 * 4.
	}
	return 0
}
func (d Dot) Weight() float64 { return 0.125 } // Assumes 8 ticks worth one normal note.

// Mods may change the duration of chart.
// Todo: implement actual calculating chart difficulties
func (c Chart) Difficulties() []float64 {
	if len(c.Notes) == 0 {
		return make([]float64, 0)
	}
	const sectionDuration = 800
	sectionCount := c.Duration()/sectionDuration + 1
	ds := make([]float64, sectionCount)
	var (
		i int
		d float64
	)
	for _, n := range c.Notes {
		for n.Time >= int64(i+1)*sectionDuration {
			ds[i] = d
			d = 0
			i++
		}
		d += n.Weight()
	}
	i, d = 0, 0

	for _, n := range c.Dots {
		for n.Time >= int64(i+1)*sectionDuration {
			ds[i] = d
			d = 0
			i++
		}
		d += n.Weight()
	}
	i, d = 0, 0

	for _, n := range c.Shakes {
		for n.Time >= int64(i+1)*sectionDuration {
			ds[i] = d
			d = 0
			i++
		}
		// Gives uniform difficulty for Shake note.
		t := int64(i) * sectionDuration
		start := n.Time // Lower bound to the section in time.
		if start < t {
			start = t
		}
		end := n.Time + n.Duration // Upper bound to the section in time.
		if end > t+sectionDuration {
			end = t + sectionDuration
		}
		var rate float64
		if n.Duration > 0 {
			rate = float64(end-start) / float64(n.Duration)
		}
		d += n.Weight() * rate
	}
	i, d = 0, 0
	return ds
}
