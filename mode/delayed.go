package mode

import (
	"math"
	"time"
)

const (
	DelayedModeExp = iota // aka DelayedModeConverge
	DelayedModeLinear
)

type Delayed struct {
	Mode     int
	Feedback float64
	Source   float64
	Delayed  float64

	countdown    int
	maxCountdown int
}

func NewDelayed() Delayed {
	const transDuration = 400 * time.Millisecond
	return Delayed{
		Mode:         DelayedModeExp,
		Feedback:     0,
		countdown:    ToTick(transDuration),
		maxCountdown: ToTick(transDuration),
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
