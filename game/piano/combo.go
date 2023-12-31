package piano

import "github.com/hndada/gosu/game"

func NewComboOptions(stage StageOptions) game.ComboOptions {
	opts := game.ComboOptions{
		Scale:    0.75,
		X:        stage.X,
		DigitGap: -1,
		Y:        0.40,
		Persist:  false,
		Bounce:   0.85,
	}
	return opts
}
