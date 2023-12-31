package game

type Judgment struct {
	Window int32
	Weight float64
}

type Judgments struct {
	Judgments  []Judgment
	blank      int
	miss       int
	missWindow int32
	Counts     []int
}

// blank is preferred to be -1,
// so that windows of Judgments are in order.
func NewJudgments(js []Judgment) Judgments {
	miss := len(js) - 1
	return Judgments{
		Judgments:  js,
		blank:      -1,
		miss:       miss,
		missWindow: js[miss].Window,
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
func (js Judgments) Judge(e int32, at KeyActionType) int {
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

func (js Judgments) Evaluate(e int32) int {
	if e < 0 {
		e *= -1
	}
	for i, j := range js.Judgments {
		if e <= j.Window {
			return i
		}
	}
	return js.blank
}
