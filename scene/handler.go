package scene

import (
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
)

var (
	VolumeMusic ctrl.KeyHandler
	VolumeSound ctrl.KeyHandler
	Brightness  ctrl.KeyHandler
	Offset      ctrl.KeyHandler
	SpeedScales []ctrl.KeyHandler
)

func init() {
	VolumeMusic = ctrl.KeyHandler{
		Handler: ctrl.FloatHandler{
			Value: &mode.S.VolumeMusic,
			Min:   0,
			Max:   1,
			Unit:  0.05,
		},
		Modifiers: []input.Key{input.KeyControlLeft},
		Keys:      [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight},
		Sounds:    [2]audios.Sounder{UserSkin.Toggle[0], UserSkin.Toggle[1]},
		Volume:    &mode.S.VolumeSound,
	}
	VolumeSound = ctrl.KeyHandler{
		Handler: ctrl.FloatHandler{
			Value: &mode.S.VolumeSound,
			Min:   0,
			Max:   1,
			Unit:  0.05,
		},
		Modifiers: []input.Key{input.KeyAltLeft},
		Keys:      [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight},
		Sounds:    [2]audios.Sounder{UserSkin.Toggle[0], UserSkin.Toggle[1]},
		Volume:    &mode.S.VolumeSound,
	}
	Brightness = ctrl.KeyHandler{
		Handler: ctrl.FloatHandler{
			Value: &mode.S.BackgroundBrightness,
			Min:   0,
			Max:   1,
			Unit:  0.1,
		},
		Modifiers: []input.Key{input.KeyControlLeft},
		Keys:      [2]input.Key{input.KeyO, input.KeyP},
		Sounds:    [2]audios.Sounder{UserSkin.Toggle[0], UserSkin.Toggle[1]},
		Volume:    &mode.S.VolumeSound,
	}
	Offset = ctrl.KeyHandler{
		Handler: ctrl.Int64Handler{
			Value: &mode.S.Offset,
			Min:   -300,
			Max:   300,
			Loop:  false,
		},
		Modifiers: []input.Key{input.KeyShiftLeft},
		Keys:      [2]input.Key{input.KeyArrowLeft, input.KeyArrowRight},
		Sounds:    [2]audios.Sounder{UserSkin.Tap, UserSkin.Tap},
		Volume:    &mode.S.VolumeSound,
	}
	SpeedScales = make([]ctrl.KeyHandler, 2)
	for i, v := range []*float64{&piano.S.SpeedScale, &drum.S.SpeedScale} {
		SpeedScales[i] = ctrl.KeyHandler{
			Handler: ctrl.FloatHandler{
				Value: v,
				Min:   0.1,
				Max:   2,
				Unit:  0.1,
			},
			Modifiers: []input.Key{},
			Keys:      [2]input.Key{input.KeyPageDown, input.KeyPageUp},
			Sounds: [2]audios.Sounder{
				UserSkin.Transition[0], UserSkin.Transition[1]},
			Volume: &mode.S.VolumeSound,
		}
	}
}
