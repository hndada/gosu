package draws

// Time is a point of time.
// Duration a length of time.
func ToTick(time int64, tps int) int { return int(float64(time) / 1000 * float64(tps)) }
func ToTime(tick int, tps int) int64 { return int64(float64(tick) / float64(tps) * 1000) }

// MaxTick, Period == {+, +}: finite drawing with animation. e.g., Judgment.
// MaxTick, Period == {+, 0}: finite drawing with no animation. e.g., Combo.
// MaxTick, Period == {0, +}: infinite drawing with animation. e.g., Dancer.
// MaxTick, Period == {0, 0}: infinite drawing with no animation. e.g., Stage.
type Timer struct {
	Tick    int
	MaxTick int
	Period  int
}

func NewTimer(maxTick, period int) Timer {
	return Timer{
		Tick:    maxTick,
		MaxTick: maxTick,
		Period:  period,
	}
}

func (t *Timer) Ticker() {
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
func (t *Timer) Reset() { t.Tick = 0 }

// For visual effects.
func (t Timer) Age() float64 {
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
func (t Timer) IsDone() bool { return t.MaxTick != 0 && t.Tick == t.MaxTick }

// For Animation.
func (t Timer) Frame(animation Animation) Sprite {
	if len(animation) == 0 {
		return Sprite{}
	}
	if t.Period == 0 {
		return animation[0]
	}
	progress := float64(t.Tick%t.Period) / float64(t.Period)
	count := float64(len(animation))
	return animation[int(progress*count)]
}
