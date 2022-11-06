package draws

// MaxTick, Period == {+, +}: finite drawing with animation. e.g., Judgment.
// MaxTick, Period == {+, 0}: finite drawing with no animation. e.g., Combo.
// MaxTick, Period == {0, +}: infinite drawing with animation. e.g., Dancer.
// MaxTick, Period == {0, 0}: infinite drawing with no animation. e.g., Stage.
type Timer2 struct {
	Tick    int
	MaxTick int
	Period  int
}

func NewTimer2(tick, period int) Timer2 {
	return Timer2{
		Tick:    tick,
		MaxTick: tick,
		Period:  period,
	}
}

func (t *Timer2) Ticker() {
	if t.MaxTick == 0 {
		if t.Period == 0 {
			return
		}
		t.Tick++
		if t.Tick > t.Period {
			t.Tick %= t.Period
		}
		return
	}
	if t.Tick < t.MaxTick {
		t.Tick++
	}
}

// For visual effects.
func (t Timer2) Age() float64 {
	if t.MaxTick == 0 {
		return 0
	}
	return float64(t.Tick) / float64(t.MaxTick)
}
func (t Timer) Progress(start, end float64) float64 {
	if end-start == 0 {
		return 0
	}
	return float64(t.Age()-start) / float64(end-start)
}

// func (t Timer) Regress(start, end float64) float64 { return 1 - t.Progress(start, end) }

// For Animation.
func (t Timer2) Frame(sprites []Sprite) Sprite {
	if len(sprites) == 0 {
		return Sprite{}
	}
	if t.Period == 0 {
		return sprites[0]
	}
	progress := float64(t.Tick%t.Period) / float64(t.Period)
	count := float64(len(sprites))
	return sprites[int(progress*count)]
}
