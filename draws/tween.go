package draws

import "math"

// Tween calculates intermediate values between two values over a specified duration.
// Yoyo is nearly no use when each tweens is not continuous.
// Todo: Add yoyo mode?
type Tween struct {
	units   []tween
	index   int
	loop    int
	maxLoop int
	// yoyo     bool // for yoyo
	// backward bool // for yoyo
}

func NewTween(begin, change float64, maxTick int, easing TweenFunc) (tw Tween) {
	tw.AppendTween(begin, change, maxTick, easing)
	return
}

func (tw *Tween) AppendTween(begin, change float64, maxTick int, easing TweenFunc) {
	tw.units = append(tw.units, tween{
		begin:   begin,
		change:  change,
		maxTick: maxTick,
		easing:  easing,
	})
}

func (tw *Tween) SetLoop(maxLoop int) { tw.maxLoop = maxLoop }

func (tw Tween) Unit() tween { return tw.units[tw.index] }

func (tw Tween) Current() float64 {
	if len(tw.units) == 0 {
		return 0
	}
	return tw.Unit().Current()
}

func (tw *Tween) Tick() {
	if tw.IsFinished() {
		return
	}
	if len(tw.units) == 0 {
		return
	}

	// Process the current Tween
	tw.units[tw.index].Tick()
	if !tw.units[tw.index].IsFinished() {
		return
	}

	// Process the next Tween
	if tw.index < len(tw.units)-1 {
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

func (tw *Tween) Reset() {
	for i := range tw.units {
		tw.units[i].tick = 0
	}
	tw.index = 0
	tw.loop = 0
}

type tween struct {
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

func (tw *tween) Tick() {
	if tw.IsFinished() {
		return
	}
	if tw.tick < tw.maxTick {
		tw.tick++
	}
}

func (tw tween) IsFinished() bool { return tw.tick >= tw.maxTick }

func (tw tween) Current() float64 {
	return tw.easing(tw.tick, tw.begin, tw.change, tw.maxTick)
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
