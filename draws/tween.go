package draws

import "math"

// Tween encapsulates the easing function along with timing data.
// It is used to animate between two values over a specified duration.
// Todo: Add yoyo mode?
type Tween struct {
	tick    int
	start   float64
	change  float64 // end - start
	maxTick int
	easing  TweenFunc
	// backward bool // for yoyo
}

// Easing function requires 4 arguments:
// current time (tick), start and change values, and duration (max tick).
type TweenFunc func(tick int, start, change float64, maxTick int) float64

func NewTween(start, change float64, maxTick int, easing TweenFunc) Tween {
	return Tween{
		start:   start,
		change:  change,
		maxTick: maxTick,
		easing:  easing,
	}
}

// IsFinished returns false if the loop is infinite.
func (tw Tween) IsFinished() bool { return tw.tick >= tw.maxTick }

// func (tw *Tween) SetChange(change float64) {
// 	tw.start = tw.start + tw.change
// 	tw.change = change
// 	tw.Reset()
// }

// func (tw *Tween) Reset() {
// 	tw.tick = 0
// 	tw.change = 0
// }

func (tw *Tween) Tick() {
	if tw.IsFinished() {
		return
	}

	if tw.tick < tw.maxTick {
		tw.tick++
	}
	tw.change = tw.easing(tw.tick, tw.start, tw.change, tw.maxTick)
}

// Tweens has separate loop variables apart from each Tween
// because loop can be either each Tween's loop or whole Tweens' loop.

// Yoyo is nearly no use when each tweens is not continuous.
// Hence, yoyo is implemented in Tweens only, not in Tween.
// Loop is also implemented in Tweens only for readability.
type Tweens struct {
	Tweens []Tween
	index  int

	loop    int
	maxLoop int
	// yoyo     bool // for yoyo
	// backward bool // for yoyo
}

func NewTweens(tws ...Tween) Tweens {
	return Tweens{Tweens: tws}
}

func (tws *Tweens) SetLoop(maxLoop int, yoyo bool) {
	tws.maxLoop = maxLoop
	// tws.yoyo = yoyo
}

// IsFinished returns false if the loop is infinite.
func (tws Tweens) IsFinished() bool {
	return tws.maxLoop != 0 && tws.loop >= tws.maxLoop
}

// func (tws *Tweens) Reset() {
// 	for i := range tws.Tweens {
// 		tws.Tweens[i].Reset()
// 	}
// 	tws.index = 0
// 	tws.loop = 0
// }

func (tws *Tweens) Tick() {
	if tws.IsFinished() {
		return
	}
	if len(tws.Tweens) == 0 {
		return
	}

	// Process the current Tween
	tws.Tweens[tws.index].Tick()
	if !tws.Tweens[tws.index].IsFinished() {
		return
	}

	if tws.index < len(tws.Tweens)-1 {
		tws.index++
	} else {
		tws.loop++
		if tws.loop < tws.maxLoop {
			tws.index = 0
		}
	}

	// // Move to the next Tween if the current one is finished
	// switch {
	// case !tws.yoyo:
	// 	// Standard behavior: increment tick until maxTick
	// 	if tws.index < len(tws.Tweens)-1 {
	// 		tws.index++
	// 	} else {
	// 		tws.loop++
	// 		if tws.loop < tws.maxLoop {
	// 			tws.index = 0
	// 		}
	// 	}
	// case tws.yoyo && !tws.backward:
	// 	// Yoyo mode - increasing tick
	// 	if tws.index < len(tws.Tweens)-1 {
	// 		tws.index++
	// 	} else {
	// 		tws.backward = true
	// 	}
	// case tws.yoyo && tws.backward:
	// 	// Yoyo mode - decreasing tick
	// 	if tws.index > 0 {
	// 		tws.index--
	// 	} else {
	// 		tws.backward = false
	// 		tws.loop++
	// 	}
	// }
}

// Easing functions
func EaseLinear(start, end, current float64, maxTick int) float64 {
	return (end - start) / float64(maxTick)
}

// y = start + (end-start) * (1-exp(-k/x))
func EaseOutExponential(tick int, start, change float64, maxTick int) float64 {
	// Decayed.go
	// k := math.Log(math.Abs(end - start)) // steepness
	// factor := 1 - math.Exp(-k/float64(maxTick))
	// return (end - start) * factor

	if tick >= maxTick {
		return start + change
	}
	k := 10.0 // Constant for exponential decay rate
	return start + change*(-math.Pow(2, -k*float64(tick)/float64(maxTick))+1)

}
