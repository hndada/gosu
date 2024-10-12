package game

import (
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/plays"
	"github.com/hndada/gosu/ui"
)

// Array is never shallow copied.
var LeftRightControls = [2]ui.Control{
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

var DownUpControls = [2]ui.Control{
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
	BackgroundBrightness ui.KeyNumberHandler[float32]
	DebugPrint           ui.KeyBoolHandler

	Mode        ui.KeyNumberHandler[int]
	SubMode     ui.KeyNumberHandler[int]
	SpeedScales []ui.KeyNumberHandler[float64]
}

func NewHandlers(opts *Options, kbs *ui.KeyboardState) *Handlers {
	return &Handlers{
		MusicVolume:          newMusicVolumeHandler(opts, kbs),
		SoundVolumeScale:     newSoundVolumeScaleHandler(opts, kbs),
		MusicOffset:          newMusicOffsetHandler(opts, kbs),
		BackgroundBrightness: newBackgroundBrightnessHandler(opts, kbs),
		DebugPrint:           newDebugPrintHandler(opts, kbs),

		Mode:        newModeHandler(opts, kbs),
		SubMode:     newSubModeHandlers(opts, kbs)[0],
		SpeedScales: newSpeedScaleHandlers(opts, kbs),
	}
}

// Control contains sound filename.
// SoundPlayer contains a pointer to sound volume scale.
func newMusicVolumeHandler(opts *Options, kbs *ui.KeyboardState) ui.KeyNumberHandler[float64] {
	return ui.KeyNumberHandler[float64]{
		NumberController: ui.NumberController[float64]{
			Value: &opts.MusicVolume,
			Min:   0,
			Max:   1,
			Unit:  0.05,
		},
		KeyListener: *ui.NewKeyListener(
			kbs,
			[]input.Key{input.KeyControlLeft},
			DownUpControls[:],
		),
	}
}

func newSoundVolumeScaleHandler(opts *Options, kbs *ui.KeyboardState) ui.KeyNumberHandler[float64] {
	return ui.KeyNumberHandler[float64]{
		NumberController: ui.NumberController[float64]{
			Value: &opts.SoundVolumeScale,
			Min:   0,
			Max:   1,
			Unit:  0.05,
		},
		KeyListener: *ui.NewKeyListener(
			kbs,
			[]input.Key{input.KeyAltLeft},
			DownUpControls[:],
		),
	}
}

func newMusicOffsetHandler(opts *Options, kbs *ui.KeyboardState) ui.KeyNumberHandler[int32] {
	ctrls := LeftRightControls
	ctrls[ui.Decrease].SoundFilename = SoundTransitionDown
	ctrls[ui.Increase].SoundFilename = SoundTransitionUp

	return ui.KeyNumberHandler[int32]{
		NumberController: ui.NumberController[int32]{
			Value: &opts.MusicOffset,
			Min:   -200,
			Max:   200,
			Unit:  1,
		},
		KeyListener: *ui.NewKeyListener(
			kbs,
			[]input.Key{input.KeyShiftLeft},
			ctrls[:],
		),
	}
}

func newBackgroundBrightnessHandler(opts *Options, kbs *ui.KeyboardState) ui.KeyNumberHandler[float32] {
	return ui.KeyNumberHandler[float32]{
		NumberController: ui.NumberController[float32]{
			Value: &opts.BackgroundBrightness,
			Min:   0,
			Max:   1,
			Unit:  0.1,
		},
		KeyListener: *ui.NewKeyListener(
			kbs,
			[]input.Key{input.KeyTab},
			LeftRightControls[:],
		),
	}
}

func newDebugPrintHandler(opts *Options, kbs *ui.KeyboardState) ui.KeyBoolHandler {
	ctrl := ui.Control{
		Key:           input.KeyF12,
		Type:          ui.Toggle,
		SoundFilename: SoundTaps,
	}
	return ui.KeyBoolHandler{
		BoolController: ui.BoolController{
			Value: &opts.DebugPrint,
		},
		KeyListener: *ui.NewKeyListener(
			kbs,
			[]input.Key{},
			[]ui.Control{ctrl},
		),
	}
}

// Todo: loop = true?
func newModeHandler(opts *Options, kbs *ui.KeyboardState) ui.KeyNumberHandler[int] {
	ctrl := ui.Control{
		Key:           input.KeyF1,
		Type:          ui.Toggle,
		SoundFilename: SoundSwipes,
	}
	return ui.KeyNumberHandler[int]{
		NumberController: ui.NumberController[int]{
			Value: &opts.Mode,
			Min:   plays.ModeAll,
			Max:   plays.ModePiano,
			Unit:  1,
		},
		KeyListener: *ui.NewKeyListener(
			kbs,
			[]input.Key{},
			[]ui.Control{ctrl},
		),
	}
}

// Todo: loop = true?
func newSubModeHandlers(opts *Options, kbs *ui.KeyboardState) []ui.KeyNumberHandler[int] {
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
	for mode := plays.ModePiano; mode <= plays.ModePiano; mode++ {
		min, max := 0, 0
		switch mode {
		case plays.ModePiano:
			min, max = 4, 10
		}

		hs = append(hs, ui.KeyNumberHandler[int]{
			NumberController: ui.NumberController[int]{
				Value: &opts.SubMode,
				Min:   min,
				Max:   max,
				Unit:  1,
			},
			KeyListener: *ui.NewKeyListener(
				kbs,
				[]input.Key{},
				ctrls[:],
			),
		})
	}
	return hs
}

func newSpeedScaleHandlers(opts *Options, kbs *ui.KeyboardState) []ui.KeyNumberHandler[float64] {
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
	for mode := plays.ModePiano; mode <= plays.ModePiano; mode++ {
		var ptr *float64
		switch mode {
		case plays.ModePiano:
			ptr = &opts.Piano.SpeedScale
		}
		hs = append(hs, ui.KeyNumberHandler[float64]{
			NumberController: ui.NumberController[float64]{
				Value: ptr,
				Min:   0.5,
				Max:   2.5,
				Unit:  0.1,
			},
			KeyListener: *ui.NewKeyListener(
				kbs,
				[]input.Key{},
				ctrls[:],
			),
		})
	}
	return hs
}
