package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type HitLightsRes struct {
	frames draws.Frames
}

func (br *HitLightsRes) Load(fsys fs.FS) {
	fname := "piano/lighting/hit.png"
	br.frames = draws.NewFramesFromFile(fsys, fname)
}

type HitLightsOpts struct {
	Scale   float64
	kx      []float64
	y       float64
	Opacity float32
}

func NewHitLightsOpts(keys KeysOpts) HitLightsOpts {
	return HitLightsOpts{
		Scale:   1.0,
		kx:      keys.kx,
		y:       keys.y,
		Opacity: 0.5,
	}
}

type HitLightsComp struct {
	keysAnim []draws.Animation
}

func NewHitLightsComp(res HitLightsRes, opts HitLightsOpts) (comp HitLightsComp) {
	keyCount := len(opts.kx)
	comp.keysAnim = make([]draws.Animation, keyCount)
	for k := range comp.keysAnim {
		a := draws.NewAnimation(res.frames, 150)
		a.MultiplyScale(opts.Scale)
		a.Locate(opts.kx[k], opts.y, draws.CenterBottom) // -HintHeight
		a.ColorScale.Scale(1, 1, 1, opts.Opacity)
		a.SetLoop(1)
		comp.keysAnim[k] = a
	}
	return
}

// Tail also makes hit lighting on.
func (comp *HitLightsComp) Update(kji []int) {
	for k, ji := range kji {
		if ji < miss {
			comp.keysAnim[k].Reset()
		}
	}
}

// HitLightsComp.Draw draws hit lights when Normal is Hit or Tail is Released.
func (comp HitLightsComp) Draw(dst draws.Image) {
	for _, a := range comp.keysAnim {
		if a.IsFinished() {
			continue
		}
		a.Draw(dst)
	}
}
