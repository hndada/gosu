package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

type HitLightsRes struct {
	frames draws.Frames
}

func (br *HitLightsRes) Load(fsys fs.FS) {
	fname := "piano/lighting/hit"
	br.frames = draws.NewFramesFromFilename(fsys, fname)
}

type HitLightsOpts struct {
	Scale   float64
	xs      []float64
	y       float64
	Opacity float64
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
		anim := draws.NewAnimation(res.frames, mode.ToTick(150))
		anim.MultiplyScale(opts.Scale)
		anim.Locate(opts.xs[k], opts.y, draws.CenterBottom) // -HintHeight
		for i := range anim.Sprites {
			anim.Sprites[i].Color.Scale(1, 1, 1, opts.Opacity)
		}
		comp.anims[k] = anim
	}
	return
}
