package game

// Judge returns index of judgment.
// Judge judges in normal style: Whether a player hits a key in time.
// Late hit makes negative time error.
func Judge(js []Judgment, e int32, a KeyActionType) int {
	miss := len(js) - 1
	blank := len(js)
	switch {
	case e > js[miss].Window:
		return blank
	case e < -js[miss].Window:
		return miss
	default: // in range
		if a == Hit {
			return Evaluate(js, e)
		} else {
			return blank
		}
	}
}

func Evaluate(js []Judgment, e int32) int {
	blank := len(js)
	if e < 0 {
		e *= -1
	}
	for i, j := range js {
		if e <= j.Window {
			return i
		}
	}
	// Returns blank when the input is out of widest range.
	return blank
}
