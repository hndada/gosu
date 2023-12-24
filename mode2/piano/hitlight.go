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
	xs      []float64
	y       float64
	Opacity float32
}

func NewHitLightsOpts(keys KeysOpts) HitLightsOpts {
	return HitLightsOpts{
		Scale:   1.0,
		xs:      keys.xs,
		y:       keys.BaselineY,
		Opacity: 0.5,
	}
}

type HitLightsComp struct {
	anims []draws.Animation
}

func NewHitLightsComp(res HitLightsRes, opts HitLightsOpts) (comp HitLightsComp) {
	keyCount := len(opts.xs)
	comp.anims = make([]draws.Animation, keyCount)
	for k := range comp.anims {
		a := draws.NewAnimation(res.frames, 150)
		a.MultiplyScale(opts.Scale)
		a.Locate(opts.xs[k], opts.y, draws.CenterBottom) // -HintHeight
		a.ColorScale.Scale(1, 1, 1, opts.Opacity)
		a.SetLoop(1)
		comp.anims[k] = a
	}
	return
}

func (comp *HitLightsComp) Update(keyOns []bool) {
	for k, ok := range keyOns {
		if ok {
			comp.anims[k].Reset()
		}
	}
}

// HitLightsComp.Draw draws hit lights when Normal is Hit or Tail is Release.
func (comp HitLightsComp) Draw(dst draws.Image) {
	for _, a := range comp.anims {
		if a.IsFinished() {
			continue
		}
		a.Draw(dst)
	}
}
