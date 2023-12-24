package piano

import mode "github.com/hndada/gosu/mode2"

func NewComboOpts(keys KeysOpts) mode.ComboOpts {
	opts := mode.ComboOpts{
		Scale:    0.75,
		X:        keys.StageX,
		Y:        keys.BaselineY,
		DigitGap: -1,
		Bounce:   0.85,
		Persist:  false,
	}
	return opts
}
