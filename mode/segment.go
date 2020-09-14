package mode

import "math"

type Segment struct {
	xMin      float64
	xMax      float64
	slope     float64
	intercept float64
}

func GetSegments(xPoints, yPoints []float64) []Segment {
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

func SolveY(segments []Segment, x float64) float64 {
	if x < 0 {
		panic("negative x")
		// panic(&ValError{"Input X to segments", Ftoa(x), ErrSyntax})
	}
	for _, segment := range segments {
		if x > segment.xMax || x < segment.xMin {
			continue
		}
		return segment.intercept + segment.slope*x
	}
	panic("cannot reach with given x")
	// panic(&ValError{"Input X to segments", Ftoa(x), ErrSyntax})
}

func SolveX(segments []Segment, y float64) []float64 {
	var x float64
	xValues := make([]float64, 0, 1)
scan:
	for _, segment := range segments {
		switch segment.slope {
		case 0:
			if y == segment.intercept {
				for _, v := range xValues {
					// if math.Abs(v-segment.xMin) < 0.0001 {
					if v == segment.xMin {
						continue scan // already same value was put
					}
				}
				xValues = append(xValues, segment.xMin)
			}
		default:
			x = (y - segment.intercept) / segment.slope
			x = Round(x, 2)
			if x >= segment.xMin && x < segment.xMax { // check interval
				xValues = append(xValues, x)
			}
		}
	}
	return xValues
}

func Round(v float64, decimal int) float64 {
	scale := math.Pow(10, float64(decimal))
	return math.Round(v*scale) / scale
}
