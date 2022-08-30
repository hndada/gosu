package piano

import "github.com/hndada/gosu"

// Todo: Variate factors based on difficulty-skewed charts
var (
	FlowScoreFactor     float64 = 0.5 // a
	AccScoreFactor      float64 = 5   // b
	KoolRateScoreFactor float64 = 2   // c
)

// Mods may change the duration of chart.
// Todo: implement actual calculating chart difficulties
func (c Chart) Difficulties() []float64 {
	if len(c.Notes) == 0 {
		return make([]float64, 0)
	}
	endTime := c.Notes[len(c.Notes)-1].Time
	ds := make([]float64, 0, endTime/gosu.SliceDuration+1)
	t := c.Notes[0].Time
	var d float64
	for _, n := range c.Notes {
		for n.Time > t+gosu.SliceDuration {
			ds = append(ds, d)
			d = 0
			t += gosu.SliceDuration
		}
		switch n.Type {
		case gosu.Tail:
			d += 0.15
		default:
			d += 1
		}
	}
	return ds
}

var FingerMap = map[int][]int{
	0:  {},
	1:  {0},
	2:  {1, 1},
	3:  {1, 0, 1},
	4:  {2, 1, 1, 2},
	5:  {2, 1, 0, 1, 2},
	6:  {3, 2, 1, 1, 2, 3},
	7:  {3, 2, 1, 0, 1, 2, 3},
	8:  {4, 3, 2, 1, 1, 2, 3, 4},
	9:  {4, 3, 2, 1, 0, 1, 2, 3, 4},
	10: {4, 3, 2, 1, 0, 0, 1, 2, 3, 4},
}

func init() {
	for k := 2; k <= 8; k++ {
		FingerMap[k|LeftScratch] = append([]int{FingerMap[k-1][0] + 1}, FingerMap[k-1]...)
		FingerMap[k|RightScratch] = append(FingerMap[k-1], FingerMap[k-1][k-2]+1)
	}
}
