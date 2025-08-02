package drum

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
			ds[i] += d
			d = 0
			i++
		}
		d += n.Weight()
		// Todo: combine to n.Weight()?
		if n.Size == Big {
			d += 0.1
		}
	}
	i, d = 0, 0

	for _, n := range c.Dots {
		for n.Time >= int64(i+1)*sectionDuration {
			ds[i] += d
			d = 0
			i++
		}
		d += n.Weight()
	}
	i, d = 0, 0

	for _, n := range c.Shakes {
		for n.Time >= int64(i+1)*sectionDuration {
			ds[i] += d
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
