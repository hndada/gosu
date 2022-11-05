package draws

const (
	Left = iota
	Center
	Right
)

const (
	Top = iota
	Middle
	Bottom
)

type Origin struct{ X, Y int }
type Align struct{ X, Y int }
