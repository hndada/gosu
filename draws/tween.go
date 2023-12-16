package draws

import "math"

const (
	TweenLoop = 0
	TweenOnce = 1
)

// Tween encapsulates the easing function along with timing data.
// It is used to animate between two values over a specified duration.
type Tween struct {
	tick     int
	maxTick  int
	loop     int
	maxLoop  int
	yoyo     bool // for yoyo
	backward bool // for yoyo

	start   float64
	current float64
	end     float64
	easing  TweenFunc
}

type TweenFunc func(tick, maxTick int, start, current, end float64) float64

func (tw *Tween) Update(end float64) {
	tw.updateTick()
	if tw.end != end {
		tw.start = tw.current
		tw.end = end
		tw.Reset()
	}
	tw.current = tw.easing(tw.tick, tw.maxTick, tw.start, tw.current, tw.end)
}

func (tw *Tween) updateTick() {
	if tw.yoyo {
		if tw.backward {
			// Decrease tick in yoyo mode
			if tw.tick > 0 {
				tw.tick--
			} else {
				// Change direction and handle loop count
				tw.backward = false
				if tw.loop < tw.maxLoop || tw.maxLoop == 0 {
					tw.loop++
				}
			}
		} else {
			// Increase tick and switch to backward at maxTick
			if tw.tick < tw.maxTick {
				tw.tick++
			} else {
				tw.backward = true
			}
		}
	} else {
		// Standard behavior: increment tick until maxTick
		if tw.tick < tw.maxTick {
			tw.tick++
		} else if tw.loop < tw.maxLoop || tw.maxLoop == 0 {
			// Reset tick and increment loop count
			tw.tick = 0
			tw.loop++
		}
	}
}

func (tw Tween) IsFinished() bool {
	return tw.loop >= tw.maxLoop && tw.tick >= tw.maxTick
}

func (tw *Tween) Reset() {
	tw.tick = 0
	tw.loop = 0
	tw.backward = false
}

// Easing functions
func EaseLinear(tick, maxTick int, start, current, end float64) float64 {
	return current + (end-start)/float64(maxTick)
}

func EaseOutExponential(tick, maxTick int, start, current, end float64) float64 {
	diff := math.Abs(end - start)
	factor := 1 - math.Exp(-math.Log(diff)/float64(maxTick))
	return current + (end-start)*factor
}
