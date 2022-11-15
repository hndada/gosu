package mode

import (
	"image/color"

	"github.com/hndada/gosu/input"
)

var (
	ColorKool = color.NRGBA{0, 170, 242, 255}   // Blue
	ColorCool = color.NRGBA{85, 251, 255, 255}  // Skyblue
	ColorGood = color.NRGBA{51, 255, 40, 255}   // Lime
	ColorBad  = color.NRGBA{244, 177, 0, 255}   // Yellow
	ColorMiss = color.NRGBA{109, 120, 134, 255} // Gray
)

type Judgment struct {
	Flow   float64
	Acc    float64
	Window int64
	// Extra  bool // For distinguishing Big note at Drum mode.
}

// Is returns whether j and j2 are equal by its window size.
func (j Judgment) Is(j2 Judgment) bool { return j.Window == j2.Window }

// Valid returns whether j is not a blank judgment by its window size.
func (j Judgment) Valid() bool { return j.Window != 0 }

// func inRange(td int64, j Judgment) bool { return td < j.Window && td > -j.Window }

// Verdict for normal notes, e.g., Note, Head at Piano mode.
func Verdict(js []Judgment, a input.KeyAction, td int64) Judgment {
	Miss := js[len(js)-1]
	switch {
	case td > Miss.Window:
		// Does nothing.
	case td < -Miss.Window:
		return Miss
	default: // In range
		if a == input.Hit {
			return Judge(js, td)
		}
	}
	return Judgment{}
}

func Judge(js []Judgment, td int64) Judgment {
	if td < 0 {
		td *= -1
	}
	for _, j := range js {
		if td <= j.Window {
			return j
		}
	}
	return Judgment{} // Returns None when the input is out of widest range
}
