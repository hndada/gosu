package draws

type Align int

const (
	Start Align = iota
	Center
	End
)

type Aligns struct{ X, Y Align }

// https://go.dev/play/p/6FsxRuznEtE
var (
	LeftTop      = Aligns{Start, Start}
	LeftMiddle   = Aligns{Start, Center}
	LeftBottom   = Aligns{Start, End}
	CenterTop    = Aligns{Center, Start}
	CenterMiddle = Aligns{Center, Center}
	CenterBottom = Aligns{Center, End}
	RightTop     = Aligns{End, Start}
	RightMiddle  = Aligns{End, Center}
	RightBottom  = Aligns{End, End}
)
