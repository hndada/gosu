package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

type HoldLightsRes struct {
	frames draws.Frames
}

func (br *HoldLightsRes) Load(fsys fs.FS) {
	fname := "piano/lighting/hold"
	br.frames = draws.NewFramesFromFilename(fsys, fname)
}

type HoldLightsOpts struct {
	Scale   float64
	xs      []float64
	y       float64
	Opacity float64
}

func NewHoldLightsOpts(keys KeysOpts) HoldLightsOpts {
	return HoldLightsOpts{
		Scale:   1.0,
		xs:      keys.xs,
		y:       keys.BaselineY,
		Opacity: 1.2,
	}
}

// field name: sprites, anims
// local name: s, a
type HoldLightsComp struct {
	anims []draws.Animation
}

func NewHoldLightsComp(res HoldLightsRes, opts HoldLightsOpts) (comp HoldLightsComp) {
	keyCount := len(opts.xs)
	comp.anims = make([]draws.Animation, keyCount)
	for k := range comp.anims {
		anim := draws.NewAnimation(res.frames, mode.ToTick(300))
		anim.MultiplyScale(opts.Scale)
		anim.Locate(opts.xs[k], opts.y, draws.CenterBottom)
		for i := range anim.Sprites {
			anim.Sprites[i].Color.Scale(1, 1, 1, opts.Opacity)
		}
		comp.anims[k] = anim
	}
	return
}
