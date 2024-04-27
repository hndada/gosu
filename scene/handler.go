package scene

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/ui"
)

// Array is never shallow copied.
var leftRightControls = [2]ui.Control{
	{
		Key:           input.KeyArrowLeft,
		Type:          ui.Decrease,
		SoundFilename: SoundToggleOff,
	},
	{
		Key:           input.KeyArrowRight,
		Type:          ui.Increase,
		SoundFilename: SoundToggleOn,
	},
}
var downUpControls = [2]ui.Control{
	{
		Key:           input.KeyArrowDown,
		Type:          ui.Decrease,
		SoundFilename: SoundToggleOff,
	},
	{
		Key:           input.KeyArrowUp,
		Type:          ui.Increase,
		SoundFilename: SoundToggleOn,
	},
}

// Controller first, then Listener.
type Handlers struct {
	MusicVolume          ui.KeyNumberHandler[float64]
	SoundVolumeScale     ui.KeyNumberHandler[float64]
	MusicOffset          ui.KeyNumberHandler[int32]
	BackgroundBrightness ui.KeyNumberHandler[float64]
	DebugPrint           ui.KeyBoolHandler

	Mode        ui.KeyNumberHandler[int]
	SubMode     ui.KeyNumberHandler[int]
	SpeedScales []ui.KeyNumberHandler[float64]
}

func NewHandlers(opts Options, states States) *Handlers {
	return &Handlers{
		MusicVolume:          newMusicVolumeHandler(opts, states),
		SoundVolumeScale:     newSoundVolumeScaleHandler(opts, states),
		MusicOffset:          newMusicOffsetHandler(opts, states),
		BackgroundBrightness: newBackgroundBrightnessHandler(opts, states),
		DebugPrint:           newDebugPrintHandler(opts, states),

		Mode:        newModeHandler(opts, states),
		SubMode:     newSubModeHandlers(opts, states)[0],
		SpeedScales: newSpeedScaleHandlers(opts, states),
	}
}

// Control contains sound filename.
// SoundPlayer contains a pointer to sound volume scale.
func newMusicVolumeHandler(opts Options, states States) ui.KeyNumberHandler[float64] {
	return ui.KeyNumberHandler[float64]{
		NumberController: ui.NumberController[float64]{
			Value: &opts.Audio.MusicVolume,
			Min:   0,
			Max:   1,
			Unit:  0.05,
		},
		KeyListener: *ui.NewKeyListener(
			states.KeyboardStatus,
			[]input.Key{input.KeyControlLeft},
			downUpControls[:],
		),
	}
}

func newSoundVolumeScaleHandler(opts Options, states States) ui.KeyNumberHandler[float64] {
	return ui.KeyNumberHandler[float64]{
		NumberController: ui.NumberController[float64]{
			Value: &opts.Audio.SoundVolumeScale,
			Min:   0,
			Max:   1,
			Unit:  0.05,
		},
		KeyListener: *ui.NewKeyListener(
			states.KeyboardStatus,
			[]input.Key{input.KeyAltLeft},
			downUpControls[:],
		),
	}
}

func newMusicOffsetHandler(opts Options, states States) ui.KeyNumberHandler[int32] {
	ctrls := leftRightControls
	ctrls[ui.Decrease].SoundFilename = SoundTransitionDown
	ctrls[ui.Increase].SoundFilename = SoundTransitionUp

	return ui.KeyNumberHandler[int32]{
		NumberController: ui.NumberController[int32]{
			Value: &opts.Audio.MusicOffset,
			Min:   -200,
			Max:   200,
			Unit:  1,
		},
		KeyListener: *ui.NewKeyListener(
			states.KeyboardStatus,
			[]input.Key{input.KeyShiftLeft},
			ctrls[:],
		),
	}
}

func newBackgroundBrightnessHandler(opts Options, states States) ui.KeyNumberHandler[float64] {
	return ui.KeyNumberHandler[float64]{
		NumberController: ui.NumberController[float64]{
			Value: &opts.Screen.BackgroundBrightness,
			Min:   0,
			Max:   1,
			Unit:  0.1,
		},
		KeyListener: *ui.NewKeyListener(
			states.KeyboardStatus,
			[]input.Key{input.KeyTab},
			leftRightControls[:],
		),
	}
}

func newDebugPrintHandler(opts Options, states States) ui.KeyBoolHandler {
	ctrl := ui.Control{
		Key:           input.KeyF12,
		Type:          ui.Toggle,
		SoundFilename: SoundTaps,
	}
	return ui.KeyBoolHandler{
		BoolController: ui.BoolController{
			Value: &opts.Screen.DebugPrint,
		},
		KeyListener: *ui.NewKeyListener(
			states.KeyboardStatus,
			[]input.Key{},
			[]ui.Control{ctrl},
		),
	}
}

// Todo: loop = true?
func newModeHandler(opts Options, states States) ui.KeyNumberHandler[int] {
	ctrl := ui.Control{
		Key:           input.KeyF1,
		Type:          ui.Toggle,
		SoundFilename: SoundSwipes,
	}
	return ui.KeyNumberHandler[int]{
		NumberController: ui.NumberController[int]{
			Value: &opts.Game.Mode,
			Min:   game.ModeAll,
			Max:   game.ModePiano,
			Unit:  1,
		},
		KeyListener: *ui.NewKeyListener(
			states.KeyboardStatus,
			[]input.Key{},
			[]ui.Control{ctrl},
		),
	}
}

// Todo: loop = true?
func newSubModeHandlers(opts Options, states States) []ui.KeyNumberHandler[int] {
	// Declare ctrls as a array for better readability.
	ctrls := [2]ui.Control{
		{
			Key:           input.KeyF2,
			Type:          ui.Decrease,
			SoundFilename: SoundSwipes,
		},
		{
			Key:           input.KeyF3,
			Type:          ui.Increase,
			SoundFilename: SoundSwipes,
		},
	}

	hs := make([]ui.KeyNumberHandler[int], 0, 3)
	for mode := game.ModePiano; mode <= game.ModePiano; mode++ {
		min, max := 0, 0
		switch mode {
		case game.ModePiano:
			min, max = 4, 10
		}

		hs = append(hs, ui.KeyNumberHandler[int]{
			NumberController: ui.NumberController[int]{
				Value: &opts.Game.SubMode,
				Min:   min,
				Max:   max,
				Unit:  1,
			},
			KeyListener: *ui.NewKeyListener(
				states.KeyboardStatus,
				[]input.Key{},
				ctrls[:],
			),
		})
	}
	return hs
}

func newSpeedScaleHandlers(opts Options, states States) []ui.KeyNumberHandler[float64] {
	// Declare ctrls as a array for better readability.
	ctrls := [2]ui.Control{
		{
			Key:           input.KeyPageDown,
			Type:          ui.Decrease,
			SoundFilename: SoundToggleOff,
		},
		{
			Key:           input.KeyPageUp,
			Type:          ui.Increase,
			SoundFilename: SoundToggleOn,
		},
	}

	hs := make([]ui.KeyNumberHandler[float64], 0, 3)
	for mode := game.ModePiano; mode <= game.ModePiano; mode++ {
		var ptr *float64
		switch mode {
		case game.ModePiano:
			ptr = &opts.Game.Piano.SpeedScale
		}
		hs = append(hs, ui.KeyNumberHandler[float64]{
			NumberController: ui.NumberController[float64]{
				Value: ptr,
				Min:   0.5,
				Max:   2.5,
				Unit:  0.1,
			},
			KeyListener: *ui.NewKeyListener(
				states.KeyboardStatus,
				[]input.Key{},
				ctrls[:],
			),
		})
	}
	return hs
}
