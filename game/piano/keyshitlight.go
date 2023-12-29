package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type KeysHitLightResources struct {
	frames draws.Frames
}

func (br *KeysHitLightResources) Load(fsys fs.FS) {
	fname := "piano/lighting/hit.png"
	br.frames = draws.NewFramesFromFile(fsys, fname)
}

type KeysHitLightOptions struct {
	Scale   float64
	kx      []float64
	y       float64
	Opacity float32
}

func NewKeysHitLightOptions(keys KeysOptions) KeysHitLightOptions {
	return KeysHitLightOptions{
		Scale:   1.0,
		kx:      keys.kx,
		y:       keys.y,
		Opacity: 0.5,
	}
}

type KeysHitLightComponent struct {
	keysAnim []draws.Animation
}

func NewKeysHitLightComponent(res KeysHitLightResources, opts KeysHitLightOptions) (cmp KeysHitLightComponent) {
	keyCount := len(opts.kx)
	cmp.keysAnim = make([]draws.Animation, keyCount)
	for k := range cmp.keysAnim {
		a := draws.NewAnimation(res.frames, 150)
		a.MultiplyScale(opts.Scale)
		a.Locate(opts.kx[k], opts.y, draws.CenterBottom) // -HintHeight
		a.ColorScale.Scale(1, 1, 1, opts.Opacity)
		a.SetLoop(1)
		cmp.keysAnim[k] = a
	}
	return
}

// Tail also makes hit lighting on.
func (cmp *KeysHitLightComponent) Update(kji []int) {
	for k, ji := range kji {
		if ji < miss {
			cmp.keysAnim[k].Reset()
		}
	}
}

// KeysHitLightComponent.Draw draws hit lights when Normal is Hit or Tail is Released.
func (cmp KeysHitLightComponent) Draw(dst draws.Image) {
	for _, a := range cmp.keysAnim {
		if a.IsFinished() {
			continue
		}
		a.Draw(dst)
	}
}
