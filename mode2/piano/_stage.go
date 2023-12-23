package piano

import mode "github.com/hndada/gosu/mode2"

// Letting user set stage width directly may be more convenient.
type StageOpts struct {
	Ws map[int]float64
	w  float64
	X  float64
}

func NewStageOpts() StageOpts {
	return StageOpts{
		Ws: map[int]float64{
			1:  240,
			2:  260,
			3:  280,
			4:  300,
			5:  320,
			6:  340,
			7:  360,
			8:  380,
			9:  400,
			10: 420,
		},
		X: 0.50 * mode.ScreenW,
	}
}
