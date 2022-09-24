package drum

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
)

var SpeedKeyHandler ctrl.KeyHandler

func LoadHandlers() {
	{
		h := &ctrl.FloatHandler{
			Value:  &SpeedScale,
			Unit:   0.1,
			Min:    0.1,
			Max:    2,
			Sounds: [2][]byte{},
		}
		SpeedKeyHandler = ctrl.KeyHandler{
			Handler: h,
		}
		SpeedKeyHandler.SetKeys(
			[]ebiten.Key{ebiten.KeyControlLeft},
			[2]ebiten.Key{ebiten.KeyPageDown, ebiten.KeyPageUp},
		)
	}
}
