package draws

type Timer2 interface {
	Ticker()
	Age() float64
	Progress(start, end float64) float64
	IsDone() bool
	Reset()
}

// tick and maxTick contains more meanings than standalone countdown.
// countdown cannot tell the zero value that is consumed or set.
type baseTimer struct {
	tick    int
	maxTick int
}

func newBaseTimer(maxTick int) baseTimer {
	return baseTimer{maxTick: maxTick}
}

// For visual effects.
func (t baseTimer) Age() float64 {
	if t.maxTick == 0 {
		return 0
	}
	return float64(t.tick) / float64(t.maxTick)
}

func (t baseTimer) Progress(start, end float64) float64 {
	if end-start == 0 {
		return 0
	}
	return float64(t.Age()-start) / float64(end-start)
}

func (t *baseTimer) Reset() { t.tick = 0 }

func (t baseTimer) IsDone() bool { return t.tick >= t.maxTick }

type FiniteTimer struct {
	baseTimer
}

func NewFiniteTimer(maxTick int) *FiniteTimer {
	return &FiniteTimer{newBaseTimer(maxTick)}
}

func (t *FiniteTimer) Ticker() {
	if t.tick < t.maxTick {
		t.tick++
	}
}

type InfiniteTimer struct {
	baseTimer
}

func (t *InfiniteTimer) Ticker() {
	if t.maxTick > 0 {
		t.tick = (t.tick % t.maxTick) + 1
	}
}

// MaxTick, Period == {+, +}: finite drawing with animation. e.g., Judgment.
// MaxTick, Period == {+, 0}: finite drawing with no animation. e.g., Combo.
// MaxTick, Period == {0, +}: infinite drawing with animation. e.g., Dancer.
// MaxTick, Period == {0, 0}: infinite drawing with no animation. e.g., Stage.
