package drum

import "github.com/hndada/gosu"

// Todo: Variate factors based on difficulty-skewed charts
var (
	FlowScoreFactor     float64 = 0.5 // a
	AccScoreFactor      float64 = 5   // b
	KoolRateScoreFactor float64 = 2   // c
)

const (
	RollTickWeight  float64 = 0.125
	ShakeTickWeight float64 = 0.03125 // 0.125 * 0.25
)

// Mods may change the duration of chart.
// Todo: implement actual calculating chart difficulties
func (c Chart) Difficulties() []float64 {
	if len(c.Notes) == 0 {
		return make([]float64, 0)
	}
	endTime := c.Notes[len(c.Notes)-1].Time
	ds := make([]float64, 0, endTime/gosu.SliceDuration+1)
	t := c.Notes[0].Time
	var d float64
	for _, n := range c.Notes {
		for n.Time > t+gosu.SliceDuration {
			ds = append(ds, d)
			d = 0
			t += gosu.SliceDuration
		}
		d += n.Weight()
		if n.Type != Head && n.Type != Shake {
			continue
		}
		// Gives uniform difficulty per time.
		// start and end is to give difficulty bound to current section.
		start := n.Time
		if start < t {
			start = t
		}
		end := n.Time2
		if end > t+gosu.SliceDuration {
			end = t + gosu.SliceDuration
		}
		// Tick is proportional to BPM.
		beats := float64(end-start) * n.ScaledBPM / 60000
		switch n.Type {
		case Head:
			// One beat has 4 Roll ticks.
			// RollTickWeight = 0.125
			// Assumes 8 ticks worth one normal note.
			ticks := beats * 4
			d += ticks * RollTickWeight
		case Shake:
			// One beat has 3 Shake ticks.
			// ShakeTickWeight = 0.03125, which is way smaller than RollTickWeight.
			// Shake is apparently easier than Roll,
			// since Shake doesn't follow the beat to hit.
			ticks := beats * 3
			d += ticks * ShakeTickWeight
		}
	}
	return ds
}

func (n Note) Weight() float64 {
	switch n.Type {
	case Normal:
		if n.Big {
			return 1.1
		} else {
			return 1.0
		}
	default:
		return 0
	}
}
