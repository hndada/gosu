package draws

import "math"

// Tween encapsulates the easing function along with timing data.
// It is used to animate between two values over a specified duration.
type Tween struct {
	start   float64
	current float64
	end     float64
	tick    int
	maxTick int
	easing  TweenFunc

	loop     int
	maxLoop  int
	yoyo     bool // for yoyo
	backward bool // for yoyo
}

type TweenFunc func(start, current, end float64, tick, maxTick int) float64

func NewTween(start, end float64, maxTick int, easing TweenFunc) Tween {
	return Tween{
		start:   start,
		current: start,
		end:     end,
		tick:    0,
		maxTick: maxTick,
		easing:  easing,
	}
}

func (tw *Tween) SetLoop(maxLoop int, yoyo bool) {
	tw.maxLoop = maxLoop
	tw.yoyo = yoyo
}

func (tw *Tween) UpdateEnd(end float64) {
	tw.start = tw.current
	tw.end = end
	tw.Reset()
}

func (tw *Tween) Tick() {
	if tw.IsFinished() {
		return
	}

	switch {
	case !tw.yoyo:
		// Standard behavior: increment tick until maxTick
		if tw.tick < tw.maxTick {
			tw.tick++
		} else {
			tw.loop++
			if tw.loop < tw.maxLoop {
				tw.tick = 0
			}
		}
	case tw.yoyo && !tw.backward:
		// Yoyo mode - increasing tick
		if tw.tick < tw.maxTick {
			tw.tick++
		} else {
			tw.backward = true
		}
	case tw.yoyo && tw.backward:
		// Yoyo mode - decreasing tick
		if tw.tick > 0 {
			tw.tick--
		} else {
			tw.backward = false
			tw.loop++
		}
	}

	// Update the current value based on the easing function
	tw.current = tw.easing(tw.start, tw.current, tw.end, tw.tick, tw.maxTick)
}

// IsFinished returns false if the loop is infinite.
func (tw Tween) IsFinished() bool {
	return tw.maxLoop != 0 && tw.loop >= tw.maxLoop
}

func (tw *Tween) Reset() {
	tw.tick = 0
	tw.loop = 0
	tw.backward = false
}

// Tweens has separate loop variables apart from each Tween
// because loop can be either each Tween's loop or whole Tweens' loop.
type Tweens struct {
	Tweens []Tween
	index  int

	loop     int
	maxLoop  int
	yoyo     bool // for yoyo
	backward bool // for yoyo
}

func (tws *Tweens) SetLoop(maxLoop int, yoyo bool) {
	tws.maxLoop = maxLoop
	tws.yoyo = yoyo
}

func (tws *Tweens) Tick() {
	if tws.IsFinished() {
		return
	}
	if len(tws.Tweens) == 0 {
		return
	}

	// Process the current Tween
	tws.Tweens[tws.index].Tick()
	tw := tws.Tweens[tws.index]
	if !tw.IsFinished() {
		return
	}

	// Move to the next Tween if the current one is finished
	tws.index++
	// Handle the wrapping of the index and update the loop count
	if tws.index >= len(tws.Tweens) {
		tws.index = 0
		tws.loop++
		// In yoyo mode, reverse the direction of each Tween
		if tws.yoyo {
			for i := range tws.Tweens {
				tws.Tweens[i].backward = !tws.Tweens[i].backward
			}
		}
	}
}

func (tws Tweens) IsFinished() bool {
	return tws.maxLoop != 0 && tws.loop >= tws.maxLoop
}

func (tws *Tweens) Reset() {
	for i := range tws.Tweens {
		tws.Tweens[i].Reset()
	}
	tws.loop = 0
	tws.backward = false
}

// Easing functions
func EaseLinear(start, current, end float64, tick, maxTick int) float64 {
	return current + (end-start)/float64(maxTick)
}

func EaseOutExponential(start, current, end float64, tick, maxTick int) float64 {
	diff := math.Abs(end - start)
	factor := 1 - math.Exp(-math.Log(diff)/float64(maxTick))
	return current + (end-start)*factor
}
