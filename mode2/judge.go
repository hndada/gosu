package mode

import "github.com/hndada/gosu/input"

// Judge judges in normal style: Whether a player hits a key in time.
// Late hit makes negative time error.
func Judge(js []Judgment, e int32, a input.KeyActionType) Judgment {
	miss := js[len(js)-1]
	switch {
	case e > miss.Window:
		return blank
	case e < -miss.Window:
		return miss
	default: // in range
		if a == input.Hit {
			return Evaluate(js, e)
		} else {
			return blank
		}
	}
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
