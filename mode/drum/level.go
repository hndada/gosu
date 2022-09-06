package drum

// Todo: Variate factors based on difficulty-skewed charts
var (
	DifficultyDuration  int64   = 800
	FlowScoreFactor     float64 = 0.5 // a
	AccScoreFactor      float64 = 5   // b
	KoolRateScoreFactor float64 = 2   // c
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
	case Shake: // Shake is apparently easier than Roll, since it is free from beat.
		return 0.03125 // 0.125 * 0.25
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

	ds := make([]float64, 0, c.Duration()/DifficultyDuration+1)
	t := c.Notes[0].Time
	var d float64
	for _, n := range c.Notes {
		for n.Time > t+DifficultyDuration {
			ds = append(ds, d)
			d = 0
			t += DifficultyDuration
		}
		if n.Type != Shake {
			d += n.Weight()
			continue
		}

		// Gives uniform difficulty for Shake ticks.
		// start and end is to give difficulty bound to current section.
		start := n.Time
		if start < t {
			start = t
		}
		end := n.Time + n.Duration
		if end > t+DifficultyDuration {
			end = t + DifficultyDuration
		}
		var rate float64 = 0.0
		if end-start > 0 {
			rate = float64(n.Duration) / float64(end-start)
		}
		ticks := float64(n.Tick) * rate
		d += ticks * n.Weight()
		// beats := float64(end-start) * n.ScaledBPM / 60000
		// switch n.Type {
		// case Head:
		// 	ticks := beats * 4
		// 	d += ticks * DotWeight
		// case Shake:
		// 	ticks := beats * 3
		// 	d += ticks * ShakeWeight
		// }
	}
	ds2 := make([]float64, 0, c.Duration()/DifficultyDuration+1)
	for _, dot := range c.Dots {
		for dot.Time > t+DifficultyDuration {
			ds = append(ds, d)
			d = 0
			t += DifficultyDuration
		}
		d += dot.Weight()
	}
	for i, d2 := range ds2 {
		ds[i] += d2
	}
	return ds
}
