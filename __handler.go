package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
)

//	func NewVolumeHandler(vol *float64, keys []ebiten.Key) ctrl.F64Handler {
//		play := func() { Sounds.Play("default-hover") }
//		return ctrl.F64Handler{
//			Handler: ctrl.Handler{
//				Keys:       keys,
//				PlaySounds: []func(){play, play},
//				HoldKey:    -1,
//			},
//			Min:    0,
//			Max:    1,
//			Unit:   0.05,
//			Target: vol,
//		}
//	}
var (
	// MusicVolumeHandler    *ctrl.FloatHandler // Todo: make it unexported?
	MusicVolumeKeyHandler  ctrl.KeyHandler
	EffectVolumeKeyHandler ctrl.KeyHandler
	// SpeedKeyHandlers       []ctrl.KeyHandler
	// modeHandler    ctrl.Handler
	ModeKeyHandler ctrl.KeyHandler
)

func LoadHandlers(modeProps []ModeProp) {
	{
		h := &ctrl.FloatHandler{
			Value:  &MusicVolume,
			Unit:   0.05,
			Min:    0,
			Max:    1,
			Sounds: TransitionSounds,
		}
		MusicVolumeKeyHandler = ctrl.KeyHandler{
			Handler: h,
		}
		MusicVolumeKeyHandler.SetKeys(
			[]ebiten.Key{ebiten.KeyAlt},
			[2]ebiten.Key{ebiten.KeyArrowDown, ebiten.KeyArrowUp},
		)
	}
	{
		h := &ctrl.FloatHandler{
			Value:  &EffectVolume,
			Unit:   0.05,
			Min:    0,
			Max:    1,
			Sounds: TransitionSounds,
		}
		EffectVolumeKeyHandler = ctrl.KeyHandler{
			Handler: h,
		}
		EffectVolumeKeyHandler.SetKeys(
			[]ebiten.Key{ebiten.KeyControlLeft},
			[2]ebiten.Key{ebiten.KeyArrowDown, ebiten.KeyArrowUp},
		)
	}
	{
		h := &ctrl.IntHandler{
			Value:  &CurrentMode,
			Unit:   1,
			Min:    0,
			Max:    len(modeProps),
			Loop:   true,
			Sounds: TransitionSounds,
		}
		ModeKeyHandler = ctrl.KeyHandler{
			Handler: h,
		}
		ModeKeyHandler.SetKeys(
			[]ebiten.Key{ebiten.KeyControlLeft, ebiten.KeyAltLeft, ebiten.KeyShiftLeft},
			[2]ebiten.Key{ebiten.KeyArrowLeft, ebiten.KeyArrowRight},
		)
	}
}

// func NewSpeedHandler(speedScale *float64) ctrl.F64Handler {
// 	play := func() { Sounds.Play("default-hover") }
// 	return ctrl.F64Handler{
// 		Handler: ctrl.Handler{
// 			Keys:       []ebiten.Key{ebiten.Key9, ebiten.Key8},
// 			PlaySounds: []func(){play, play},
// 			HoldKey:    -1,
// 		},
// 		Min:    0.1,
// 		Max:    2,
// 		Unit:   0.1,
// 		Target: speedScale,
// 	}
// }

// func NewModeHandler(cursor *int, len int) ctrl.IntHandler {
// 	play := func() { Sounds.Play("default-hover") }
// 	return ctrl.IntHandler{
// 		Handler: ctrl.Handler{
// 			Keys:       []ebiten.Key{ebiten.Key0},
// 			PlaySounds: []func(){play},
// 			HoldKey:    -1,
// 		},
// 		Min:    0,
// 		Max:    len,
// 		Unit:   1,
// 		Target: cursor,
// 		Loop:   true,
// 	}
// }

// func NewVsyncSwitchHandler(b *bool) ctrl.BoolHandler {
// 	play := func() { Sounds.Play("default-hover") }
// 	return ctrl.BoolHandler{
// 		Handler: ctrl.Handler{
// 			Keys:       []ebiten.Key{ebiten.Key5},
// 			PlaySounds: []func(){play},
// 			HoldKey:    -1,
// 		},
// 		Target: b,
// 	}
// }

func NewSpeedKeyHandler(speedScale *float64) (h ctrl.KeyHandler) {
	h = ctrl.KeyHandler{
		Handler: &ctrl.FloatHandler{
			Value:  speedScale,
			Unit:   0.1,
			Min:    0.1,
			Max:    2,
			Sounds: [2][]byte{},
		},
	}
	h.SetKeys(
		[]ebiten.Key{ebiten.KeyControlLeft},
		[2]ebiten.Key{ebiten.KeyPageDown, ebiten.KeyPageUp},
	)
	return
}
