package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type KeysHoldLightResources struct {
	frames draws.Frames
}

func (br *KeysHoldLightResources) Load(fsys fs.FS) {
	fname := "piano/lighting/hold.png"
	br.frames = draws.NewFramesFromFile(fsys, fname)
}

type KeysHoldLightOptions struct {
	Scale   float64
	kx      []float64
	y       float64
	Opacity float32
}

func NewKeysHoldLightOptions(keys KeysOptions) KeysHoldLightOptions {
	return KeysHoldLightOptions{
		Scale:   1.0,
		kx:      keys.kx,
		y:       keys.y,
		Opacity: 1.2,
	}
}

type KeysHoldLightComponent struct {
	anims               []draws.Animation
	keysLongNoteHolding []bool
}

func NewKeysHoldLightComponent(res KeysHoldLightResources, opts KeysHoldLightOptions) (cmp KeysHoldLightComponent) {
	keyCount := len(opts.kx)
	cmp.anims = make([]draws.Animation, keyCount)
	for k := range cmp.anims {
		a := draws.NewAnimation(res.frames, 300)
		a.MultiplyScale(opts.Scale)
		a.Locate(opts.kx[k], opts.y, draws.CenterBottom)
		a.ColorScale.Scale(1, 1, 1, opts.Opacity)
		cmp.anims[k] = a
	}
	return
}

func (cmp *KeysHoldLightComponent) Update(kn []Note, kh []bool, ka game.KeyboardAction) {
	kln := make([]bool, len(kn))
	for k, n := range kn {
		if n.valid && n.Type == Tail {
			kln[k] = true
		}
	}
	klnh := make([]bool, len(kh))
	for k, n := range kln {
		if n && kh[k] {
			klnh[k] = true
		}
	}

	keysOld := cmp.keysLongNoteHolding
	keysNew := klnh
	for k := range klnh {
		if !keysOld[k] && keysNew[k] {
			cmp.anims[k].Reset()
		}
	}
	cmp.keysLongNoteHolding = klnh
}

func (cmp KeysHoldLightComponent) Draw(dst draws.Image) {
	for _, a := range cmp.anims {
		a.Draw(dst)
	}
}
