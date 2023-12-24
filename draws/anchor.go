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

type Anchor struct{ X, Y int }

// https://go.dev/play/p/6FsxRuznEtE
var (
	LeftTop      = Anchor{Left, Top}
	LeftMiddle   = Anchor{Left, Middle}
	LeftBottom   = Anchor{Left, Bottom}
	CenterTop    = Anchor{Center, Top}
	CenterMiddle = Anchor{Center, Middle}
	CenterBottom = Anchor{Center, Bottom}
	RightTop     = Anchor{Right, Top}
	RightMiddle  = Anchor{Right, Middle}
	RightBottom  = Anchor{Right, Bottom}
)
