package game

type Judgment struct {
	Window int32
	Weight float64
}

// Is returns whether two Judgments are equal.
// func (j Judgment) Is(j2 Judgment) bool { return j.Window == j2.Window }
// func (j Judgment) IsBlank() bool       { return j.Window == 0 }

// var blank = Judgment{}

// // the ideal number of Judgments is: 3 + 1
// const (
// 	Kool = iota
// 	Cool
// 	Good
// 	Miss // Its window is used for judging too early hit.
// )
