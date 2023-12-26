package piano

import base "github.com/hndada/gosu/game"

func NewComboOpts(keys KeysOpts) base.ComboOpts {
	opts := base.ComboOpts{
		Scale:    0.75,
		X:        keys.StageX,
		Y:        keys.BaselineY,
		DigitGap: -1,
		Bounce:   0.85,
		Persist:  false,
	}
	return opts
}
