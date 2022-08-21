package mode

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
)

func NewVolumeHandler(vol *float64) ctrl.F64Handler {
	// b, err := audios.NewBytes("skin/default-hover.wav")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// play := audios.Context.NewPlayerFromBytes(b).Play
	play := func() { Sounds.Play("default-hover") }
	return ctrl.F64Handler{
		Handler: ctrl.Handler{
			Keys:       []ebiten.Key{ebiten.Key2, ebiten.Key1},
			PlaySounds: []func(){play, play},
			HoldKey:    -1,
		},
		Min:    0,
		Max:    1,
		Unit:   0.05,
		Target: vol,
	}
}

func NewSpeedHandler(speedBase *float64) ctrl.F64Handler {
	// b, err := audios.NewBytes("skin/default-hover.wav")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// play := audios.Context.NewPlayerFromBytes(b).Play
	play := func() { Sounds.Play("default-hover") }
	return ctrl.F64Handler{
		Handler: ctrl.Handler{
			Keys:       []ebiten.Key{ebiten.Key4, ebiten.Key3},
			PlaySounds: []func(){play, play},
			HoldKey:    -1,
		},
		Min:    0.1,
		Max:    2,
		Unit:   0.1,
		Target: speedBase,
	}
}

// Todo: should Max be *int?
func NewSelectHandler(cursor *int, len int) ctrl.IntHandler {
	// b, err := audios.NewBytes("skin/default-hover.wav")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// play := audios.Context.NewPlayerFromBytes(b).Play
	play := func() { Sounds.Play("default-hover") }
	return ctrl.IntHandler{
		Handler: ctrl.Handler{
			Keys:       []ebiten.Key{ebiten.KeyDown, ebiten.KeyUp},
			PlaySounds: []func(){play, play},
			HoldKey:    -1,
		},
		Min:    0,
		Max:    len,
		Unit:   1,
		Target: cursor,
	}
}
