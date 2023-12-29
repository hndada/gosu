package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type KeysHoldLightRes struct {
	frames draws.Frames
}

func (br *KeysHoldLightRes) Load(fsys fs.FS) {
	fname := "piano/lighting/hold.png"
	br.frames = draws.NewFramesFromFile(fsys, fname)
}

type KeysHoldLightOpts struct {
	Scale   float64
	kx      []float64
	y       float64
	Opacity float32
}

func NewKeysHoldLightOpts(keys KeysOpts) KeysHoldLightOpts {
	return KeysHoldLightOpts{
		Scale:   1.0,
		kx:      keys.kx,
		y:       keys.y,
		Opacity: 1.2,
	}
}

// field name: sprites, anims
// local name: s, a
type KeysHoldLightComp struct {
	anims       []draws.Animation
	keysHolding []bool
}

func NewKeysHoldLightComp(res KeysHoldLightRes, opts KeysHoldLightOpts) (comp KeysHoldLightComp) {
	keyCount := len(opts.kx)
	comp.anims = make([]draws.Animation, keyCount)
	for k := range comp.anims {
		a := draws.NewAnimation(res.frames, 300)
		a.MultiplyScale(opts.Scale)
		a.Locate(opts.kx[k], opts.y, draws.CenterBottom)
		a.ColorScale.Scale(1, 1, 1, opts.Opacity)
		comp.anims[k] = a
	}
	return
}

func (comp *KeysHoldLightComp) Update(kh []bool) {
	olds := comp.keysHolding
	for k, new := range kh {
		if new && !olds[k] {
			comp.anims[k].Reset()
		}
	}
	comp.keysHolding = kh
}

func (comp KeysHoldLightComp) Draw(dst draws.Image) {
	for _, a := range comp.anims {
		a.Draw(dst)
	}
}
