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

// https://go.dev/play/p/6FsxRuznEtE
var (
	LeftTop      = struct{ X, Y int }{Left, Top}
	LeftMiddle   = struct{ X, Y int }{Left, Middle}
	LeftBottom   = struct{ X, Y int }{Left, Bottom}
	CenterTop    = struct{ X, Y int }{Center, Top}
	CenterMiddle = struct{ X, Y int }{Center, Middle}
	CenterBottom = struct{ X, Y int }{Center, Bottom}
	RightTop     = struct{ X, Y int }{Right, Top}
	RightMiddle  = struct{ X, Y int }{Right, Middle}
	RightBottom  = struct{ X, Y int }{Right, Bottom}
)

// type Align struct{ X, Y int }
