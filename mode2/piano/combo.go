package piano

import (
	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

func NewComboOpts(stage draws.WHXY) mode.ComboOpts {
	opts := mode.ComboOpts{
		Scale:    0.75,
		RX:       0.50,
		RY:       0.40,
		DigitGap: -1,
		Bounce:   0.85,
		Persist:  false,
	}
	opts.SetStage(stage)
	return opts
}
