package piano

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
)

var (
	// MusicVolumeKeyHandler  = gosu.MusicVolumeKeyHandler
	// EffectVolumeKeyHandler = gosu.EffectVolumeKeyHandler
	SpeedKeyHandler        ctrl.KeyHandler
	Piano4CursorKeyHandler ctrl.KeyHandler
	Piano7CursorKeyHandler ctrl.KeyHandler
)

func LoadHandler() {
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
	{
		h := &ctrl.IntHandler{
			Value:  &SpeedScale,
			Unit:   1,
			Min:    0,
			Max:    len(ModePiano7.ChartInfos),
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
