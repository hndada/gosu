package scene

import (
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/input"
)

var LeftRightKeys = [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight}
var UpDownKeys = [2]input.Key{input.KeyArrowUp, input.KeyArrowDown}

func IsEnterJustPressed() bool {
	return input.IsKeyJustPressed(input.KeyEnter) || input.IsKeyJustPressed(input.KeyNumpadEnter)
}
func IsEscapeJustPressed() bool {
	return input.IsKeyJustPressed(input.KeyEscape)
}

func NewMusicVolumeKeyHandler(cfg *Config, asset *Asset) func() bool {
	handler := ctrl.KeyHandler{
		Handler: ctrl.FloatHandler{
			Value: &cfg.MusicVolume,
			Min:   0,
			Max:   1,
			Unit:  0.05,
		},
		Modifier: input.KeyControlLeft,
		Keys:     LeftRightKeys,
		Sounds:   asset.ToggleSounds,
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewSoundVolumeKeyHandler(cfg *Config, asset *Asset) func() bool {
	handler := ctrl.KeyHandler{
		Handler: ctrl.FloatHandler{
			Value: &cfg.SoundVolume,
			Min:   0,
			Max:   1,
			Unit:  0.05,
		},
		Modifier: input.KeyAltLeft,
		Keys:     LeftRightKeys,
		Sounds:   asset.ToggleSounds,
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewMusicOffsetKeyHandler(cfg *Config, asset *Asset) func() bool {
	handler := ctrl.KeyHandler{
		Handler: ctrl.Int32Handler{
			Value: &cfg.MusicOffset,
			Min:   -200,
			Max:   200,
			Loop:  false,
			Unit:  1,
		},
		Modifier: input.KeyShiftLeft,
		Keys:     LeftRightKeys,
		Sounds:   asset.TransitionSounds,
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewBackgroundBrightnessKeyHandler(cfg *Config, asset *Asset) func() bool {
	handler := ctrl.KeyHandler{
		Handler: ctrl.FloatHandler{
			Value: &cfg.BackgroundBrightness,
			Min:   0,
			Max:   1,
			Unit:  0.1,
		},
		Modifier: input.KeyControlLeft,
		Keys:     [2]input.Key{input.KeyO, input.KeyP},
		Sounds:   asset.ToggleSounds,
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewDebugPrintKeyHandler(cfg *Config, asset *Asset) func() bool {
	handler := ctrl.KeyHandler{
		Handler: ctrl.BoolHandler{
			Value: &cfg.DebugPrint,
		},
		Modifier: input.KeyNone,
		Keys:     [2]input.Key{input.KeyF12, input.KeyF12},
		Sounds:   [2]audios.SoundPlayer{asset.TapSoundPod, asset.TapSoundPod},
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewSpeedScaleKeyHandler(cfg *Config, asset *Asset, mode int) func() bool {
	var ptr *float64
	switch mode {
	case ModePiano:
		ptr = &cfg.PianoConfig.SpeedScale
	}

	handler := ctrl.KeyHandler{
		Handler: ctrl.FloatHandler{
			Value: ptr,
			Min:   0.5,
			Max:   2.5,
			Unit:  0.1,
		},
		Modifier: input.KeyNone,
		Keys:     [2]input.Key{input.KeyPageDown, input.KeyPageUp},
		Sounds:   asset.ToggleSounds,
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewModeKeyHandler(cfg *Config, asset *Asset) func() bool {
	handler := ctrl.KeyHandler{
		Handler: ctrl.IntHandler{
			Value: &cfg.Mode,
			Min:   -1, // modeAll
			Max:   0,  // modePiano only so far
			Loop:  true,
		},
		Modifier: input.KeyNone,
		Keys:     [2]input.Key{input.KeyNone, input.KeyF1},
		Sounds:   [2]audios.SoundPlayer{asset.SwipeSoundPod, asset.SwipeSoundPod},
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewSubModeKeyHandler(cfg *Config, asset *Asset, mode int) func() bool {
	min, max := 0, 0
	switch cfg.Mode {
	case ModePiano:
		min, max = 4, 10
	}
	handler := ctrl.KeyHandler{
		Handler: ctrl.IntHandler{
			Value: &cfg.SubMode,
			Min:   min,
			Max:   max,
			Loop:  true,
		},
		Modifier: input.KeyNone,
		Keys:     [2]input.Key{input.KeyF2, input.KeyF3},
		Sounds:   [2]audios.SoundPlayer{asset.SwipeSoundPod, asset.SwipeSoundPod},
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}
