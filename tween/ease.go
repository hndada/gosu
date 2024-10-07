package tween

import (
	"math"
	"time"
)

// Easing requires 4 arguments:
// current time (t), begin (b), change (c), and duration (d).
type Easing func(t time.Duration, b, c float64, d time.Duration) float64

func EaseLinear(t time.Duration, b, c float64, d time.Duration) float64 {
	r := t.Seconds() / d.Seconds()
	return b + r*c
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
	dx := t.Seconds() / d.Seconds()
	return b + c*(1-math.Exp(-k*dx))
}
