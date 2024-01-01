package scene

import (
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

func IsEnterJustPressed() bool {
	return input.IsKeyJustPressed(input.KeyEnter) || input.IsKeyJustPressed(input.KeyNumpadEnter)
}
func IsEscapeJustPressed() bool {
	return input.IsKeyJustPressed(input.KeyEscape)
}

func NewMusicVolumeKeyHandler(cfg *Config, asset *Asset) func() bool {
	handler := ctrl.KeyHandler{
		Handler:  ctrl.NewValueHandler[float64](&cfg.MusicVolume, 0, 1, 0.05),
		Modifier: input.KeyControlLeft,
		Keys:     ctrl.KeysLeftRight,
		Sounds:   asset.ToggleSounds,
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewSoundVolumeKeyHandler(cfg *Config, asset *Asset) func() bool {
	handler := ctrl.KeyHandler{
		Handler:  ctrl.NewValueHandler[float64](&cfg.SoundVolume, 0, 1, 0.05),
		Modifier: input.KeyAltLeft,
		Keys:     ctrl.KeysLeftRight,
		Sounds:   asset.ToggleSounds,
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewMusicOffsetKeyHandler(cfg *Config, asset *Asset) func() bool {
	handler := ctrl.KeyHandler{
		Handler:  ctrl.NewValueHandler[int](&cfg.MusicOffset, -200, 200, 1),
		Modifier: input.KeyShiftLeft,
		Keys:     ctrl.KeysLeftRight,
		Sounds:   asset.TransitionSounds,
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewBackgroundBrightnessKeyHandler(cfg *Config, asset *Asset) func() bool {
	handler := ctrl.KeyHandler{
		Handler:  ctrl.NewValueHandler[float64](&cfg.BackgroundBrightness, 0, 1, 0.1),
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

func NewSpeedScaleKeyHandler(cfg *Config, asset *Asset, _mode int) func() bool {
	var ptr *float64
	switch _mode {
	case mode.ModePiano:
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
			Min:   mode.ModeAll,   // modeAll
			Max:   mode.ModePiano, // modePiano only so far
			Loop:  true,
		},
		Modifier: input.KeyNone,
		Keys:     [2]input.Key{input.KeyNone, input.KeyF1},
		Sounds:   [2]audios.SoundPlayer{asset.SwipeSoundPod, asset.SwipeSoundPod},
		Volume:   &cfg.SoundVolume,
	}
	return handler.Handle
}

func NewSubModeKeyHandler(cfg *Config, asset *Asset) func() bool {
	min, max := 0, 0
	switch cfg.Mode {
	case mode.ModePiano:
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
