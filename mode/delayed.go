package mode

import "math"

const (
	DelayedModeExp = iota // aka DelayedModeConverge
	DelayedModeLinear
)

// For displaying score with scrolling effect.
type Delayed struct {
	Mode     int
	Feedback float64

	src     *float64
	Target  float64
	Delayed float64

	countdown    int
	maxCountdown int
}

func NewDelayed(src *float64) Delayed {
	const transDuration = 400
	return Delayed{
		Mode:     DelayedModeExp,
		Feedback: 0,

		src:     src,
		Target:  *src,
		Delayed: 0, // *src,

		countdown:    ToTick(transDuration),
		maxCountdown: ToTick(transDuration),
	}
}

func (d *Delayed) Update() {
	if d.Target != *d.src {
		d.Target = *d.src
		d.countdown = d.maxCountdown
		diff := d.Target - d.Delayed
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
		d.Delayed = d.Target
		return
	}
	switch d.Mode {
	case DelayedModeExp:
		d.Delayed += (d.Target - d.Delayed) * d.Feedback
	case DelayedModeLinear:
		d.Delayed += d.Feedback
	}
	d.countdown--
}
