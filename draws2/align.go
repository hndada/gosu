package draws

type Align int

const (
	Left Align = iota
	Center
	Right
)

const (
	Top Align = iota
	Middle
	Bottom
)

type Aligns struct{ X, Y Align }

// https://go.dev/play/p/6FsxRuznEtE
var (
	LeftTop      = Aligns{Left, Top}
	LeftMiddle   = Aligns{Left, Middle}
	LeftBottom   = Aligns{Left, Bottom}
	CenterTop    = Aligns{Center, Top}
	CenterMiddle = Aligns{Center, Middle}
	CenterBottom = Aligns{Center, Bottom}
	RightTop     = Aligns{Right, Top}
	RightMiddle  = Aligns{Right, Middle}
	RightBottom  = Aligns{Right, Bottom}
)
