package scene

import (
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/input"
)

var (
	currentMode int
	currentSort int
)
var (
	modeHandler    ctrl.IntHandler
	ModeKeyHandler ctrl.KeyHandler
	sortHandler    ctrl.IntHandler
	SortKeyHandler ctrl.KeyHandler

	musicVolumeHandler    ctrl.FloatHandler
	VolumeMusicKeyHandler ctrl.KeyHandler
	effectVolumeHandler   ctrl.FloatHandler
	VolumeSoundKeyHandler ctrl.KeyHandler

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
		Volume:    &VolumeSound,
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
		Volume:    &VolumeSound,
	}

	musicVolumeHandler = ctrl.FloatHandler{
		Value: &VolumeMusic,
		Min:   0,
		Max:   1,
		Unit:  0.05,
	}
	VolumeMusicKeyHandler = ctrl.KeyHandler{
		Handler:   musicVolumeHandler,
		Modifiers: []input.Key{input.KeyControlLeft},
		Keys:      [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight},
		Sounds:    [2][]byte{ToggleSounds[0], ToggleSounds[1]},
		Volume:    &VolumeSound,
	}
	effectVolumeHandler = ctrl.FloatHandler{
		Value: &VolumeSound,
		Min:   0,
		Max:   1,
		Unit:  0.05,
	}
	VolumeSoundKeyHandler = ctrl.KeyHandler{
		Handler:   effectVolumeHandler,
		Modifiers: []input.Key{input.KeyAltLeft},
		Keys:      [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight},
		Sounds:    [2][]byte{ToggleSounds[0], ToggleSounds[1]},
		Volume:    &VolumeSound,
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
		Volume:    &VolumeSound,
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
		Volume:    &VolumeSound,
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
		Volume:    &VolumeSound,
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
		Volume:    &VolumeSound,
	}
}
