package piano

import mode "github.com/hndada/gosu/mode2"

func NewComboOpts(key KeyOpts) mode.ComboOpts {
	opts := mode.ComboOpts{
		Scale:    0.75,
		X:        key.StageX,
		Y:        key.BaselineY,
		DigitGap: -1,
		Bounce:   0.85,
		Persist:  false,
	}
	return opts
}
