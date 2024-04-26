package piano

import (
	"github.com/hndada/gosu/game"
)

type StageOptions struct {
	keyCount int
	// Ws       map[int]float64
	// w        float64
	W float64
	H float64
	X float64
	y float64 // bottom
}

var defaultStageWidths = map[int]float64{
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
}

func NewStageOptions(keyCount int) StageOptions {
	opts := StageOptions{
		keyCount: keyCount,
		W:        defaultStageWidths[keyCount],
		H:        0.90 * game.ScreenH,
		X:        0.50 * game.ScreenW,
		y:        0.90 * game.ScreenH,
	}
	// opts.w = opts.Ws[opts.keyCount]
	return opts
}
