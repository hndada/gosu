package draws

// Countdown is for drawing a sprite for a while.
type BaseDrawer struct {
	Countdown    int
	MaxCountdown int // Draw permanently when value is zero.
}

func (d *BaseDrawer) Update(reloaded bool) {
	if d.Countdown > 0 {
		d.Countdown--
	}
	if reloaded {
		d.Countdown = d.MaxCountdown
	}
}

// Todo: should I handle when MaxCountdown == 0?
func (d BaseDrawer) Age() float64 {
	return 1 - (float64(d.Countdown) / float64(d.MaxCountdown))
}
