package mode

import "github.com/hndada/gosu/input"

type Judgment struct {
	Window int64
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
func Judge(js []Judgment, timeError int64, a input.KeyActionType) Judgment {
	miss := js[len(js)-1]
	switch {
	case timeError > miss.Window:
		return blank
	case timeError < -miss.Window:
		return miss
	default: // In range
		if a == input.Hit {
			return Evaluate(js, timeError)
		}
	}
	return blank
}

func Evaluate(js []Judgment, timeError int64) Judgment {
	if timeError < 0 {
		timeError *= -1
	}
	for _, j := range js {
		if timeError <= j.Window {
			return j
		}
	}
	// Returns blank when the input is out of widest range.
	return blank
}
