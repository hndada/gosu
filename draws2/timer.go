package draws

// Timer helps drawing a sprite with visual effects.
type Timer struct {
	Countdown    int
	MaxCountdown int
}

func NewTimer(cd int) Timer {
	return Timer{
		MaxCountdown: cd,
	}
}
func (t Timer) Age() float64 {
	if t.MaxCountdown == 0 {
		return 0
	}
	return 1 - (float64(t.Countdown) / float64(t.MaxCountdown))
}
