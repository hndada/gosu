package mode

import "github.com/hndada/gosu/input"

// type Scorer interface {
// 	Check()
// 	Judge()
// 	Mark()
// }

type Judgment struct {
	Window int32
	Weight float64
}

var blank = Judgment{}

// the ideal number of Judgments is: 3 + 1
const (
	Kool = iota
	Cool
	Good
	Miss // Its window is used for judging too early hit.
)

// Is returns whether two Judgments are equal.
func (j Judgment) Is(j2 Judgment) bool { return j.Window == j2.Window }
func (j Judgment) IsBlank() bool       { return j.Window == 0 }

// Judge judges in normal style: Whether a player hits a key in time.
// Late hit makes negative time error.
func Judge(js []Judgment, e int32, a input.KeyActionType) Judgment {
	miss := js[len(js)-1]
	switch {
	case e > miss.Window:
		return blank
	case e < -miss.Window:
		return miss
	default: // In range
		if a == input.Hit {
			return Evaluate(js, e)
		}
	}
	return blank
}

func Evaluate(js []Judgment, e int32) Judgment {
	if e < 0 {
		e *= -1
	}
	for _, j := range js {
		if e <= j.Window {
			return j
		}
	}
	// Returns blank when the input is out of widest range.
	return blank
}
