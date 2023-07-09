package scene

import (
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/input"
)

var (
	MusicVolumeKeyHandler          ctrl.KeyHandler
	SoundVolumeKeyHandler          ctrl.KeyHandler
	BackgroundBrightnessKeyHandler ctrl.KeyHandler
	OffsetKeyHandler               ctrl.KeyHandler
	DebugPrintKeyHandler           ctrl.KeyHandler

	SpeedScaleKeyHandlers []ctrl.KeyHandler
)

// LoadHandler should be called after asset is loaded.
func LoadHandler(speedScales []*float64) {
	var (
		ctrlKey  = []input.Key{input.KeyControlLeft}
		altKey   = []input.Key{input.KeyAltLeft}
		shfitKey = []input.Key{input.KeyShiftLeft}

		leftRightKeys = [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight}
		upDownKeys    = [2]input.Key{input.KeyArrowUp, input.KeyArrowDown}
	)
	var (
		toggleSounds     = [2]audios.Sounder{TheAsset.Toggle[0], TheAsset.Toggle[1]}
		transitionSounds = [2]audios.Sounder{TheAsset.Transition[0], TheAsset.Transition[1]}
		tapSound         = [2]audios.Sounder{TheAsset.Tap, TheAsset.Tap}
	)
	var vol = &TheSettings.SoundVolume

	MusicVolumeKeyHandler = ctrl.KeyHandler{
		Handler: ctrl.FloatHandler{
			Value: &TheSettings.MusicVolume,
			Min:   0,
			Max:   1,
			Unit:  0.05,
		},
		Modifiers: ctrlKey,
		Keys:      leftRightKeys,
		Sounds:    toggleSounds,
		Volume:    vol,
	}
	SoundVolumeKeyHandler = ctrl.KeyHandler{
		Handler: ctrl.FloatHandler{
			Value: &TheSettings.SoundVolume,
			Min:   0,
			Max:   1,
			Unit:  0.05,
		},
		Modifiers: altKey,
		Keys:      leftRightKeys,
		Sounds:    toggleSounds,
		Volume:    vol,
	}
	BackgroundBrightnessKeyHandler = ctrl.KeyHandler{
		Handler: ctrl.FloatHandler{
			Value: &TheSettings.BackgroundBrightness,
			Min:   0,
			Max:   1,
			Unit:  0.1,
		},
		Modifiers: ctrlKey,
		Keys:      [2]input.Key{input.KeyO, input.KeyP},
		Sounds:    toggleSounds,
		Volume:    vol,
	}
	OffsetKeyHandler = ctrl.KeyHandler{
		Handler: ctrl.Int64Handler{
			Value: &TheSettings.Offset,
			Min:   -300,
			Max:   300,
			Loop:  false,
			Unit:  1,
		},
		Modifiers: shfitKey,
		Keys:      leftRightKeys,
		Sounds:    transitionSounds,
		Volume:    vol,
	}
	DebugPrintKeyHandler = ctrl.KeyHandler{
		Handler: ctrl.BoolHandler{
			Value: &TheSettings.DebugPrint,
		},
		Keys:   [2]input.Key{input.KeyF12, input.KeyF12},
		Sounds: tapSound,
		Volume: vol,
	}
	SpeedScaleKeyHandlers = make([]ctrl.KeyHandler, 0, len(speedScales))
	for _, v := range speedScales {
		handler := ctrl.KeyHandler{
			Handler: ctrl.FloatHandler{
				Value: v,
				Min:   0.5,
				Max:   2,
				Unit:  0.05,
			},
			Modifiers: ctrlKey,
			Keys:      upDownKeys,
			Sounds:    toggleSounds,
			Volume:    vol,
		}
		SpeedScaleKeyHandlers = append(SpeedScaleKeyHandlers, handler)
	}
}
