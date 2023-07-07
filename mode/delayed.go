package mode

import "math"

const TransDuration = 400 // In milliseconds.

const (
	DelayedModeExp = iota // aka DelayedModeConverge
	DelayedModeLinear
	// DelayedModeDiverge
)

type Delayed struct {
	Delayed float64
	Source  float64
	// source       *float64
	Mode         int // Todo: custom type?
	Feedback     float64
	countdown    int
	maxCountdown int
}

func NewDelayed(tps int) Delayed {
	return Delayed{
		// Delayed:      *src,
		// source:       src,
		Mode:         DelayedModeExp,
		Feedback:     0,
		countdown:    ToTick(400, tps),
		maxCountdown: ToTick(400, tps),
	}
}

func (d *Delayed) Update(src float64) {
	if d.Source != src {
		d.Source = src
		d.countdown = d.maxCountdown
		diff := d.Source - d.Delayed
		if diff < 0.1 {
			return
		}
		// The formula is to make speed of transition constant regardless of TPS.
		switch d.Mode {
		case DelayedModeExp:
			d.Feedback = 1 - math.Exp(-math.Log(diff)/float64(d.maxCountdown))
		case DelayedModeLinear:
			d.Feedback = diff / float64(d.maxCountdown)
		}
	}
	if d.countdown == 0 {
		d.Delayed = d.Source
		return
	}
	switch d.Mode {
	case DelayedModeExp:
		d.Delayed += (d.Source - d.Delayed) * d.Feedback
	case DelayedModeLinear:
		d.Delayed += d.Feedback
	}
	d.countdown--
}
