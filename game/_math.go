package game

import (
	"math"
)

func BellCurve(x float64) float64 {
	const (
		a = 1.0 // height of curve's peak
		b = 0   // position of the peak
		c = 0.1 // standart deviation controlling width of the curve
		//( lower abstract value of c -> "longer" curve)
	)
	return a * math.Exp(-math.Pow(x-b, 2)/(2.0*math.Pow(c, 2)))
}
