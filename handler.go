package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
)

var (
	currentMode int
	currentSort int

	MusicVolume          float64 = 0.25
	EffectVolume         float64 = 0.25
	BackgroundBrightness float64 = 0.6

	Offset int = -65
)
var (
	modeHandler    ctrl.IntHandler
	ModeKeyHandler ctrl.KeyHandler
	sortHandler    ctrl.IntHandler
	SortKeyHandler ctrl.KeyHandler

	musicVolumeHandler     ctrl.FloatHandler
	MusicVolumeKeyHandler  ctrl.KeyHandler
	effectVolumeHandler    ctrl.FloatHandler
	EffectVolumeKeyHandler ctrl.KeyHandler

	brightHandler    ctrl.FloatHandler
	BrightKeyHandler ctrl.KeyHandler

	offsetHandler    ctrl.IntHandler
	OffsetKeyHandler ctrl.KeyHandler
)
var (
	speedScaleHandlers   []ctrl.FloatHandler
	SpeedScaleKeyHandler ctrl.KeyHandler
)

const (
	SortByName = iota
	SortByLevel
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
		Keys:      [2]ebiten.Key{-1, ebiten.KeyF1},
		Sounds:    [2][]byte{SwipeSound, SwipeSound},
		Volume:    &EffectVolume,
	}
	sortHandler = ctrl.IntHandler{
		Value: &currentSort,
		Min:   0,
		Max:   1,
		Loop:  true,
	}
	SortKeyHandler = ctrl.KeyHandler{
		Handler:   sortHandler,
		Modifiers: []ebiten.Key{},
		Keys:      [2]ebiten.Key{-1, ebiten.KeyF2},
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
		Modifiers: []ebiten.Key{ebiten.KeyControlLeft},
		Keys:      [2]ebiten.Key{ebiten.KeyLeft, ebiten.KeyRight},
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
		Modifiers: []ebiten.Key{ebiten.KeyAltLeft},
		Keys:      [2]ebiten.Key{ebiten.KeyLeft, ebiten.KeyRight},
		Sounds:    [2][]byte{ToggleSounds[0], ToggleSounds[1]},
		Volume:    &EffectVolume,
	}

	brightHandler = ctrl.FloatHandler{
		Value: &BackgroundBrightness,
		Min:   0,
		Max:   1,
		Unit:  0.1,
	}
	BrightKeyHandler = ctrl.KeyHandler{
		Handler:   brightHandler,
		Modifiers: []ebiten.Key{ebiten.KeyControlLeft},
		Keys:      [2]ebiten.Key{ebiten.KeyO, ebiten.KeyP},
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
		Keys:      [2]ebiten.Key{ebiten.KeyPageDown, ebiten.KeyPageUp},
		Sounds:    [2][]byte{TransitionSounds[0], TransitionSounds[1]},
		Volume:    &EffectVolume,
	}

	offsetHandler = ctrl.IntHandler{
		Value: &Offset,
		Min:   -300,
		Max:   300,
		Loop:  false,
	}
	OffsetKeyHandler = ctrl.KeyHandler{
		Handler:   offsetHandler,
		Modifiers: []ebiten.Key{ebiten.KeyShiftLeft},
		Keys:      [2]ebiten.Key{ebiten.KeyLeft, ebiten.KeyRight},
		Sounds:    [2][]byte{TapSound, TapSound},
		Volume:    &EffectVolume,
	}
}
