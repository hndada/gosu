package ctrl

import (
	"math"
)

const (
	DelayedModeExp = iota // aka DelayedModeConverge
	DelayedModeLinear
	// DelayedModeDiverge // I suppose this mode would not be beloved.
)

type Delayed struct {
	Delayed   float64
	Source    float64
	Mode      int // Todo: custom type?
	Feedback  float64
	Countdown int
}

// The formula is to make speed of transition constant regardless of TPS.
func (d *Delayed) Set(v float64) {
	if d.Source != v {
		d.Countdown = transCountdown
	}
	d.Source = v
	diff := d.Source - d.Delayed
	if diff < 1e-2 {
		return
	}
	switch d.Mode {
	case DelayedModeExp:
		d.Feedback = 1 - math.Exp(-math.Log(diff)/float64(transCountdown))
	case DelayedModeLinear:
		d.Feedback = diff / float64(transCountdown)
	}
}
func (d *Delayed) Update() {
	if d.Countdown == 0 {
		d.Delayed = d.Source
		return
	}
	switch d.Mode {
	case DelayedModeExp:
		d.Delayed += (d.Source - d.Delayed) * d.Feedback
	case DelayedModeLinear:
		d.Delayed += d.Feedback
	}
	d.Countdown--
}
