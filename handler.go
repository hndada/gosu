package gosu

import (
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/input"
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
var (
	tailExtraTimeHandler    ctrl.FloatHandler
	TailExtraTimeKeyHandler ctrl.KeyHandler
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
		Modifiers: []input.Key{},
		Keys:      [2]input.Key{-1, input.KeyF1},
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
		Modifiers: []input.Key{},
		Keys:      [2]input.Key{-1, input.KeyF2},
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
		Modifiers: []input.Key{input.KeyControlLeft},
		Keys:      [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight},
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
		Modifiers: []input.Key{input.KeyAltLeft},
		Keys:      [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight},
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
		Modifiers: []input.Key{input.KeyControlLeft},
		Keys:      [2]input.Key{input.KeyO, input.KeyP},
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
		Modifiers: []input.Key{},
		Keys:      [2]input.Key{input.KeyPageDown, input.KeyPageUp},
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
		Modifiers: []input.Key{input.KeyShiftLeft},
		Keys:      [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight},
		Sounds:    [2][]byte{TapSound, TapSound},
		Volume:    &EffectVolume,
	}

	tailExtraTimeHandler = ctrl.FloatHandler{
		Value: props[ModePiano4].Settings["TailExtraTime"],
		Min:   -200,
		Max:   200,
		Unit:  10,
	}
	TailExtraTimeKeyHandler = ctrl.KeyHandler{
		Handler:   tailExtraTimeHandler,
		Modifiers: []input.Key{input.KeyF3}, // Todo: configurate by Button in future
		Keys:      [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight},
		Sounds:    [2][]byte{TapSound, TapSound},
		Volume:    &EffectVolume,
	}
}
