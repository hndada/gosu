package draws

import "fmt"

// Timer helps drawing a sprite with visual effects.
type Timer struct {
	Countdown    int
	MaxCountdown int
	Duration     int
}

func NewTimer(max, duration int) Timer {
	return Timer{
		MaxCountdown: max,
		Duration:     duration,
	}
}
func (t Timer) Age() float64 {
	if t.MaxCountdown == 0 {
		return 0
	}
	return 1 - (float64(t.Countdown) / float64(t.MaxCountdown))
}

// Works as Animation.
func (t Timer) Frame(sprites []Sprite) Sprite {
	if t.Duration == 0 {
		return sprites[0]
	}
	tick := t.MaxCountdown - t.Countdown
	rate := float64(tick%t.Duration) / float64(t.Duration)
	count := float64(len(sprites))
	fmt.Println(rate)
	return sprites[int(rate*count)]
}
