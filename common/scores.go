package common

type Scores struct {
	Score       float64
	Combo       int
	HP          float64
	Judgments   []Judgment
	JudgeCounts []int
}

func (s *Scores) CountCombo(j Judgment) {
	if j.ComboBreak {
		s.Combo = 0
	} else {
		s.Combo++
	}
}

type Judgment struct {
	Value      float64
	Penalty    float64
	HP         float64
	Window     int64
	ComboBreak bool
}

func (s Scores) Judge(t NoteType, a KeyActionState, d int64) Judgment {
	if d < 0 { // d: time difference; absolute value
		d *= -1
	}
	for _, j := range s.Judgments {
		if d <= j.Window {
			return j
		}
	}
	return Judgment{} // When a player hit too early
}

func (s *Scores) CountJudge(j Judgment) {
	for i, j2 := range s.Judgments {
		if j.Value == j2.Value { // Window is not same, temporarily
			s.JudgeCounts[i]++
			break
		}
	}
}

func (s *Scores) ScaleWindow(scale float64) {
	for i := range s.Judgments {
		w := float64(s.Judgments[i].Window)
		s.Judgments[i].Window = int64(w * scale)
	}
}
