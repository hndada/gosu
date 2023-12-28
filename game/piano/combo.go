package piano

import "github.com/hndada/gosu/game"

func NewComboOpts(keys KeysOpts) game.ComboOpts {
	opts := game.ComboOpts{
		Scale:    0.75,
		X:        keys.StageX,
		Y:        keys.BaselineY,
		DigitGap: -1,
		Bounce:   0.85,
		Persist:  false,
	}
	return opts
}
