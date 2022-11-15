package ctrl

import (
	"math"
)

const (
	DelayedModeExp = iota // aka DelayedModeConverge
	DelayedModeLinear
	// DelayedModeDiverge // I suppose this mode would not be beloved.
)

// Todo: make fields unexported
type Delayed struct {
	Delayed   float64
	Source    float64
	Mode      int // Todo: custom type?
	Feedback  float64
	Countdown int
}

func (d *Delayed) Value() float64 { return d.Delayed }
func (d *Delayed) Update(v float64) {
	if d.Source != v {
		d.setSource(v)
	}
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
func (d *Delayed) setSource(v float64) {
	d.Source = v
	d.Countdown = transCountdown
	diff := d.Source - d.Delayed
	if diff < 0.1 {
		return
	}
	// The formula is to make speed of transition constant regardless of TPS.
	switch d.Mode {
	case DelayedModeExp:
		d.Feedback = 1 - math.Exp(-math.Log(diff)/float64(transCountdown))
	case DelayedModeLinear:
		d.Feedback = diff / float64(transCountdown)
	}
}
