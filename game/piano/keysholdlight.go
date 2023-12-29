package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
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

// field name: sprites, anims
// local name: s, a
type KeysHoldLightComponent struct {
	anims       []draws.Animation
	keysHolding []bool
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

func (cmp *KeysHoldLightComponent) Update(kh []bool) {
	olds := cmp.keysHolding
	for k, new := range kh {
		if new && !olds[k] {
			cmp.anims[k].Reset()
		}
	}
	cmp.keysHolding = kh
}

func (cmp KeysHoldLightComponent) Draw(dst draws.Image) {
	for _, a := range cmp.anims {
		a.Draw(dst)
	}
}
