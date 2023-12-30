package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type HitLightsResources struct {
	frames draws.Frames
}

func (br *HitLightsResources) Load(fsys fs.FS) {
	fname := "piano/lighting/hit.png"
	br.frames = draws.NewFramesFromFile(fsys, fname)
}

type HitLightsOptions struct {
	Scale   float64
	kx      []float64
	y       float64
	Opacity float32
}

func NewHitLightsOptions(keys KeysOptions) HitLightsOptions {
	return HitLightsOptions{
		Scale:   1.0,
		kx:      keys.kx,
		y:       keys.y,
		Opacity: 0.5,
	}
}

type HitLightsComponent struct {
	keysAnim []draws.Animation
}

func NewHitLightsComponent(res HitLightsResources, opts HitLightsOptions) (cmp HitLightsComponent) {
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
func (cmp *HitLightsComponent) Update(kji []int) {
	for k, ji := range kji {
		if ji < miss {
			cmp.keysAnim[k].Reset()
		}
	}
}

// HitLightsComponent.Draw draws hit lights when Normal is Hit or Tail is Released.
func (cmp HitLightsComponent) Draw(dst draws.Image) {
	for _, a := range cmp.keysAnim {
		if a.IsFinished() {
			continue
		}
		a.Draw(dst)
	}
}
