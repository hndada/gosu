package ctrl

import (
	"math"
)

const (
	DelayedModeLinear = iota
	DelayedModeExp    // aka DelayedModeConverge
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
	d.Source = v
	diff := d.Source - d.Delayed
	d.Countdown = transCountdown
	switch d.Mode {
	case DelayedModeLinear:
		d.Feedback = diff / float64(transCountdown)
	case DelayedModeExp:
		d.Feedback = 1 - math.Exp(-math.Log(diff)/float64(transCountdown))
	}
}
func (d *Delayed) Update() {
	if d.Countdown == 0 {
		d.Delayed = d.Source
		return
	}
	switch d.Mode {
	case DelayedModeLinear:
		d.Delayed += d.Feedback
	case DelayedModeExp:
		d.Delayed = (d.Delayed - d.Source) * d.Feedback
	}
	d.Countdown--
}
