package draws

import "math"

// Tween calculates intermediate values between two values over a specified duration.
// Yoyo is nearly no use when each tweens is not continuous.
// Todo: Add yoyo mode?
type Tween struct {
	tweens  []tween
	index   int
	loop    int
	maxLoop int
	// yoyo     bool // for yoyo
	// backward bool // for yoyo
}

type tween struct {
	tick    int
	begin   float64
	change  float64
	maxTick int
	easing  TweenFunc
	// backward bool // for yoyo
}

func (tw *tween) Tick() {
	if tw.IsFinished() {
		return
	}
	if tw.tick < tw.maxTick {
		tw.tick++
	}
	tw.change = tw.easing(tw.tick, tw.begin, tw.change, tw.maxTick)
}

func (tw tween) IsFinished() bool { return tw.tick >= tw.maxTick }

// Easing function requires 4 arguments:
// current time (tick), begin and change values, and duration (max tick).
type TweenFunc func(tick int, begin, change float64, maxTick int) float64

func NewTween(begin, change float64, maxTick int, easing TweenFunc) (tw Tween) {
	tw.AppendTween(begin, change, maxTick, easing)
	return
}

func (tw *Tween) AppendTween(begin, change float64, maxTick int, easing TweenFunc) {
	tw.tweens = append(tw.tweens, tween{
		begin:   begin,
		change:  change,
		maxTick: maxTick,
		easing:  easing,
	})
}

func (tw *Tween) SetLoop(maxLoop int) { tw.maxLoop = maxLoop }

func (tw *Tween) Tick() {
	if tw.IsFinished() {
		return
	}
	if len(tw.tweens) == 0 {
		return
	}

	// Process the current Tween
	tw.tweens[tw.index].Tick()
	if !tw.tweens[tw.index].IsFinished() {
		return
	}

	// Process the next Tween
	if tw.index < len(tw.tweens)-1 {
		tw.index++
	} else {
		tw.loop++
		if tw.loop < tw.maxLoop {
			tw.index = 0
		}
	}
}

// IsFinished returns false if the loop is infinite.
func (tw Tween) IsFinished() bool {
	return tw.maxLoop != 0 && tw.loop >= tw.maxLoop
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
