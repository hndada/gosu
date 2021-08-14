package game

import (
	"math"
	"sort"
)

type Segment struct {
	xMin      float64
	xMax      float64
	slope     float64
	intercept float64
}

type Segments []Segment

func NewSegments(xPoints, yPoints []float64) []Segment {
	segments := make([]Segment, len(xPoints))
	for i := range segments {
		segments[i].xMin = xPoints[i]
		switch i {
		case len(xPoints) - 1:
			segments[i].xMax = math.Inf(1)
			segments[i].slope = 0
			segments[i].intercept = yPoints[i]
		default:
			segments[i].xMax = xPoints[i+1]
			segments[i].slope = (yPoints[i+1] - yPoints[i]) / (xPoints[i+1] - xPoints[i])
			segments[i].intercept = yPoints[i] - segments[i].slope*xPoints[i]
		}
	}
	return segments
}

func (ss Segments) SolveY(x float64) float64 {
	if x < 0 {
		panic("negative x")
	}
	for _, s := range ss {
		if x > s.xMax || x < s.xMin {
			continue
		}
		return s.intercept + s.slope*x
	}
	panic("cannot reach with given x")
}

func (ss Segments) SolveX(y float64) []float64 {
	var x float64
	xValues := make([]float64, 0, 1)
scan:
	for _, s := range ss {
		switch s.slope {
		case 0:
			if y == s.intercept {
				for _, v := range xValues {
					// if math.Abs(v-s.xMin) < 0.0001 {
					if v == s.xMin {
						continue scan // already same value was put
					}
				}
				xValues = append(xValues, s.xMin)
			}
		default:
			round := func(v float64, decimal int) float64 {
				scale := math.Pow(10, float64(decimal))
				return math.Round(v*scale) / scale
			}
			x = (y - s.intercept) / s.slope
			x = round(x, 2)                // 둘째자리에서 반올림
			if x >= s.xMin && x < s.xMax { // check interval
				xValues = append(xValues, x)
			}
		}
	}
	return xValues
}

// Difficulty relating
func DecayFactor(decayBase float64, time int64) float64 {
	return math.Pow(decayBase, float64(time)/1000)
}

// Maclaurin series
func WeightedSum(series []float64, weightDecay float64) float64 {
	sort.Slice(series, func(i, j int) bool { return series[i] > series[j] })
	sum, weight := 0.0, 1.0
	for _, term := range series {
		sum += weight * term
		weight *= weightDecay
	}
	return sum
}
