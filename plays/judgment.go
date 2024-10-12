package plays

type JudgmentKind int

type Judgment struct {
	Window int32
	Weight float64
}

type Judgments struct {
	Judgments  []Judgment
	blank      JudgmentKind
	miss       JudgmentKind
	missWindow int32
	Counts     []int
}

// blank is preferred to be len(js), so that blank comes after miss.
func NewJudgments(js []Judgment) Judgments {
	return Judgments{
		Judgments:  js,
		blank:      JudgmentKind(len(js)),
		miss:       JudgmentKind(len(js) - 1),
		missWindow: js[len(js)-1].Window,
		Counts:     make([]int, len(js)),
	}
}

// e stands for time error. It decreases as the time goes by.
// In other word, late hit makes negative time error.
func (js Judgments) IsTooEarly(e int32) bool { return e > js.missWindow }
func (js Judgments) IsTooLate(e int32) bool  { return e < -js.missWindow }
func (js Judgments) IsInRange(e int32) bool {
	return !js.IsTooEarly(e) && !js.IsTooLate(e)
}

// Judge returns index of judgment.
// Judge judges in normal style: Whether a player hits a key in time.
func (js Judgments) Judge(e int32, at KeyActionType) JudgmentKind {
	switch {
	case js.IsTooEarly(e):
		return js.blank
	case js.IsTooLate(e):
		return js.miss
	case js.IsInRange(e):
		if at == Hit {
			return js.Evaluate(e)
		}
	}
	return js.blank
}

func (js Judgments) Evaluate(e int32) JudgmentKind {
	if e < 0 {
		e *= -1
	}
	for kind, j := range js.Judgments {
		if e <= j.Window {
			return JudgmentKind(kind)
		}
	}
	return js.blank
}
