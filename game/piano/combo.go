package piano

import "github.com/hndada/gosu/game"

func NewComboOptions(stage StageOptions) game.ComboOptions {
	opts := game.ComboOptions{
		Scale:    0.75,
		X:        stage.X,
		Y:        0.40,
		DigitGap: -1,
		Bounce:   0.85,
		Persist:  false,
	}
	return opts
}
