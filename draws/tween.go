package draws

import (
	"math"
	"time"

	"github.com/hndada/gosu/times"
)

type tween struct {
	startTime time.Time
	begin     float64
	change    float64
	duration  time.Duration
	easing    TweenFunc
	// backward bool // for yoyo
}

// Easing function requires 4 arguments:
// current time (t), begin (b), change (c), and duration (d).
type TweenFunc func(t time.Duration, b, c float64, d time.Duration) float64

func (tw tween) current() float64 {
	return tw.easing(times.Since(tw.startTime), tw.begin, tw.change, tw.duration)
}

func (tw tween) endTime() time.Time {
	return tw.startTime.Add(tw.duration)
}

func (tw tween) isFinished() bool {
	return tw.endTime().Before(times.Now())
}

// Tween calculates intermediate values between two values over a specified duration.
// Yoyo is nearly no use when each tweens is not continuous.
type Tween struct {
	units   []tween
	index   int
	loop    int
	maxLoop int
	// yoyo     bool // for yoyo
	// backward bool // for yoyo
}

func NewTween(begin, change float64, duration time.Duration, easing TweenFunc) (tw Tween) {
	tw.Add(begin, change, duration, easing)
	return
}

func (tw Tween) endTime() time.Time {
	if len(tw.units) == 0 {
		return times.Now()
	}
	lastTween := tw.units[len(tw.units)-1]
	return lastTween.startTime.Add(lastTween.duration)
}

// AppendXxx feels like to return a struct.
func (tw *Tween) Add(begin, change float64, duration time.Duration, easing TweenFunc) {
	tw.units = append(tw.units, tween{
		startTime: tw.endTime(),
		begin:     begin,
		change:    change,
		duration:  duration,
		easing:    easing,
	})
}

func (tw *Tween) SetLoop(maxLoop int) { tw.maxLoop = maxLoop }

func (tw *Tween) Current() float64 {
	if len(tw.units) == 0 {
		return 0
	}

	for tw.units[tw.index].isFinished() {
		if tw.index < len(tw.units)-1 {
			tw.index++
		} else {
			tw.loop++
			if tw.loop < tw.maxLoop {
				tw.rewind()
			}
		}
	}

	return tw.units[tw.index].current()
}

func (tw *Tween) rewind() {
	for i := range tw.units {
		if i == 0 {
			tw.units[i].startTime = times.Now()
		} else {
			prev := tw.units[i-1]
			tw.units[i].startTime = prev.endTime()
		}
	}
	tw.index = 0
}

func (tw *Tween) Reset() {
	tw.rewind()
	tw.loop = 0
}

// IsFinished returns false if the loop is infinite.
func (tw Tween) IsFinished() bool {
	return tw.maxLoop != 0 && tw.loop >= tw.maxLoop
}

// Todo: need to be tested
func (tw *Tween) Finish() {
	tw.loop = tw.maxLoop
	tw.index = len(tw.units) - 1
}

// Easing functions
// begin + change*dx
func EaseLinear(t time.Duration, b, c float64, d time.Duration) float64 {
	dx := float64(t) / float64(d)
	return b + c*dx
}

// begin + change*(1-math.Exp(-k*dx))
func EaseOutExponential(t time.Duration, b, c float64, d time.Duration) float64 {
	if t >= d {
		return b + c
	}

	// k, steepness, is regardless of the number of steps.
	// https://go.dev/play/p/NnGiHCfPfD-
	// k := math.Log(math.Abs(change)) // delayed.go

	const k = -6.93 // exp(-6.93) ~= pow(2, -10)
	dx := float64(t) / float64(d)
	return b + c*(1-math.Exp(-k*dx))
}

func (tw *Tween) Yoyo() {
	// // Move to the next Tween if the current one is finished
	// switch {
	// case !tw.yoyo:
	// 	// Standard behavior: increment tick until maxTick
	// 	if tw.index < len(tw.Tweens)-1 {
	// 		tw.index++
	// 	} else {
	// 		tw.loop++
	// 		if tw.loop < tw.maxLoop {
	// 			tw.index = 0
	// 		}
	// 	}
	// case tw.yoyo && !tw.backward:
	// 	// Yoyo mode - increasing tick
	// 	if tw.index < len(tw.Tweens)-1 {
	// 		tw.index++
	// 	} else {
	// 		tw.backward = true
	// 	}
	// case tw.yoyo && tw.backward:
	// 	// Yoyo mode - decreasing tick
	// 	if tw.index > 0 {
	// 		tw.index--
	// 	} else {
	// 		tw.backward = false
	// 		tw.loop++
	// 	}
	// }
}
