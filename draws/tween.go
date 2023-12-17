package draws

import "math"

// Tween calculates intermediate values between two values over a specified duration.
type Tween struct {
	tick    int
	begin   float64
	change  float64
	maxTick int
	easing  TweenFunc
	// backward bool // for yoyo
}

// Easing function requires 4 arguments:
// current time (tick), begin and change values, and duration (max tick).
type TweenFunc func(tick int, begin, change float64, maxTick int) float64

func NewTween(begin, change float64, maxTick int, easing TweenFunc) Tween {
	return Tween{
		begin:   begin,
		change:  change,
		maxTick: maxTick,
		easing:  easing,
	}
}

// IsFinished returns false if the loop is infinite.
func (tw Tween) IsFinished() bool { return tw.tick >= tw.maxTick }

func (tw *Tween) Tick() {
	if tw.IsFinished() {
		return
	}
	if tw.tick < tw.maxTick {
		tw.tick++
	}
	tw.change = tw.easing(tw.tick, tw.begin, tw.change, tw.maxTick)
}

// Easing functions
// begin + change*dx
func EaseLinear(tick int, begin, change float64, maxTick int) float64 {
	dx := float64(tick) / float64(maxTick)
	return begin + change*dx
}

// begin + change*(1-math.Exp(-k*dx))
func EaseOutExponential(tick int, begin, change float64, maxTick int) float64 {
	if tick >= maxTick {
		return begin + change
	}

	// By setting k like below, the number of steps will be constant.
	k := math.Log(math.Abs(change)) // steepness
	dx := float64(tick) / float64(maxTick)
	return begin + change*(1-math.Exp(-k*dx))
}

// Yoyo is nearly no use when each tweens is not continuous.
// Hence, yoyo is implemented in Tweens only, not in Tween.
// Loop is also implemented in Tweens only for readability.
// Todo: Add yoyo mode?
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
