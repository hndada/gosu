package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
)

var (
	currentMode int

	MusicVolume  float64 = 0.25
	EffectVolume float64 = 0.25
)
var (
	modeHandler    ctrl.IntHandler
	ModeKeyHandler ctrl.KeyHandler

	musicVolumeHandler     ctrl.FloatHandler
	MusicVolumeKeyHandler  ctrl.KeyHandler
	effectVolumeHandler    ctrl.FloatHandler
	EffectVolumeKeyHandler ctrl.KeyHandler

	speedScaleHandlers   []ctrl.FloatHandler
	SpeedScaleKeyHandler ctrl.KeyHandler
)

func LoadHandlers(props []ModeProp) {
	modeHandler = ctrl.IntHandler{
		Value: &currentMode,
		Min:   0,
		Max:   len(props) - 1,
		Loop:  true,
	}
	ModeKeyHandler = ctrl.KeyHandler{
		Handler:   modeHandler,
		Modifiers: []ebiten.Key{},
		Keys:      [2]ebiten.Key{ebiten.Key0, ebiten.KeySpace},
		Sounds:    [2][]byte{SwipeSound, SwipeSound},
		Volume:    &EffectVolume,
	}

	musicVolumeHandler = ctrl.FloatHandler{
		Value: &MusicVolume,
		Min:   0,
		Max:   1,
		Unit:  0.05,
	}
	MusicVolumeKeyHandler = ctrl.KeyHandler{
		Handler:   musicVolumeHandler,
		Modifiers: []ebiten.Key{},
		Keys:      [2]ebiten.Key{ebiten.KeyQ, ebiten.KeyW},
		Sounds:    [2][]byte{ToggleSounds[0], ToggleSounds[1]},
		Volume:    &EffectVolume,
	}
	effectVolumeHandler = ctrl.FloatHandler{
		Value: &EffectVolume,
		Min:   0,
		Max:   1,
		Unit:  0.05,
	}
	EffectVolumeKeyHandler = ctrl.KeyHandler{
		Handler:   effectVolumeHandler,
		Modifiers: []ebiten.Key{},
		Keys:      [2]ebiten.Key{ebiten.KeyA, ebiten.KeyS},
		Sounds:    [2][]byte{ToggleSounds[0], ToggleSounds[1]},
		Volume:    &EffectVolume,
	}

	speedScaleHandlers = make([]ctrl.FloatHandler, len(props))
	for i, prop := range props {
		speedScaleHandlers[i] = ctrl.FloatHandler{
			Value: prop.SpeedScale,
			Min:   0.1,
			Max:   2,
			Unit:  0.1,
		}
	}
	SpeedScaleKeyHandler = ctrl.KeyHandler{
		Handler:   speedScaleHandlers[currentMode],
		Modifiers: []ebiten.Key{},
		Keys:      [2]ebiten.Key{ebiten.KeyZ, ebiten.KeyX},
		Sounds:    [2][]byte{TransitionSounds[0], TransitionSounds[1]},
		Volume:    &EffectVolume,
	}
}
