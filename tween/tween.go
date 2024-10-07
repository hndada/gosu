package tween

import (
	"time"

	"github.com/hndada/gosu/times"
)

// Tween (in-betweening) calculates a value between two values at a certain time.
// Unit is a single tweening operation.
type Unit struct {
	Begin    float64
	Change   float64
	Duration time.Duration
	Easing   Easing
}

func (u Unit) Value(elapsed time.Duration) float64 {
	return u.Easing(elapsed, u.Begin, u.Change, u.Duration)
}

type Tween struct {
	Units []Unit
	index int

	MaxLoop int // 0 means infinite looping
	loop    int

	startTime time.Time
	starts    []time.Duration
	ends      []time.Duration
}

func (tw *Tween) Add(begin, change float64, duration time.Duration, easing Easing) {
	tw.Units = append(tw.Units, Unit{begin, change, duration, easing})
}

// Start initializes the tween sequence,
// setting the startTime and resetting indices.
func (tw *Tween) Start() {
	tw.index = 0
	tw.loop = 0
	tw.startTime = time.Now()

	if len(tw.starts) == len(tw.Units) {
		return
	}
	var t time.Duration
	tw.starts = make([]time.Duration, len(tw.Units))
	tw.ends = make([]time.Duration, len(tw.Units))
	for i, u := range tw.Units {
		tw.starts[i] = t
		t += u.Duration
		tw.ends[i] = t
	}
}

// Index is the only factor that determines if the tween is finished.
func (tw *Tween) Stop()           { tw.index = len(tw.Units) }
func (tw Tween) IsFinished() bool { return tw.index >= len(tw.Units) }
func (tw *Tween) Index() int      { return tw.index }

// Update calculates the current tween value based on the elapsed time.
func (tw *Tween) Update() {
	if len(tw.starts) != len(tw.Units) {
		tw.Start()
	}

	// Loop through the units to find the active one
	elapsed := time.Since(tw.startTime)
	for tw.index < len(tw.Units) {
		if elapsed < tw.ends[tw.index] {
			// Current unit is still in progress
			return
		}

		// Move to the next unit
		tw.index++
		if tw.index < len(tw.Units) {
			continue
		}

		// If all units are done, check for looping
		if tw.MaxLoop == 0 || tw.loop < tw.MaxLoop-1 {
			tw.loop++
			tw.index = 0
			tw.startTime = tw.startTime.Add(tw.ends[len(tw.ends)-1])
			elapsed = time.Since(tw.startTime)
		}
	}
}

func (tw Tween) Value() float64 {
	if len(tw.Units) == 0 {
		return 0
	}

	if tw.index < len(tw.Units) {
		elapsed := times.Since(tw.startTime)
		unitElapsed := elapsed - tw.starts[tw.index]
		return tw.Units[tw.index].Value(unitElapsed)
	}

	last := tw.Units[len(tw.Units)-1]
	return last.Begin + last.Change
}
