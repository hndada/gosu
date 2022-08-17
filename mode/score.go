package mode

import "github.com/hndada/gosu/input"

type Judgment struct {
	Flow   float64
	Acc    float64
	Window int64
}

type Result struct {
	FinalScore float64
	FlowScore  float64
	AccScore   float64
	ExtraScore float64

	JudgmentCounts []int
	Flows          float64 // Sum of Flow
	Accs           float64
	Extras         float64 // Kool rate, for example.
	MaxCombo       int

	FlowMarks []float64 // Length is around 100 ~ 200.
}

// Verdict for normal notes: Note, Head at long note.
func Verdict(js []Judgment, a input.KeyAction, td int64) Judgment {
	Miss := js[len(js)-1]
	switch {
	case td > Miss.Window:
		// Does nothing
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
	if td < 0 { // Absolute value
		td *= -1
	}
	for _, j := range js {
		if td <= j.Window {
			return j
		}
	}
	return Judgment{} // Returns None when the input is out of widest range
}

// // MarkedNoteCount is for calculating ratio.
// func (s ScenePlay) MarkedNoteCount() int {
// 	sum := 0
// 	for _, c := range s.JudgmentCounts {
// 		sum += c
// 	}
// 	return sum
// }

// func inRange(td int64, j Judgment) bool { return td < j.Window && td > -j.Window }
