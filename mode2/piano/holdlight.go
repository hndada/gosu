package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type HoldLightsRes struct {
	frames draws.Frames
}

func (br *HoldLightsRes) Load(fsys fs.FS) {
	fname := "piano/lighting/hold.png"
	br.frames = draws.NewFramesFromFile(fsys, fname)
}

type HoldLightsOpts struct {
	Scale   float64
	xs      []float64
	y       float64
	Opacity float32
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
	anims    []draws.Animation
	keyHolds []bool
}

func NewHoldLightsComp(res HoldLightsRes, opts HoldLightsOpts) (comp HoldLightsComp) {
	keyCount := len(opts.xs)
	comp.anims = make([]draws.Animation, keyCount)
	for k := range comp.anims {
		a := draws.NewAnimation(res.frames, 300)
		a.MultiplyScale(opts.Scale)
		a.Locate(opts.xs[k], opts.y, draws.CenterBottom)
		a.ColorScale.Scale(1, 1, 1, opts.Opacity)
		comp.anims[k] = a
	}
	return
}

func (comp *HoldLightsComp) Update(keyHolds []bool) {
	olds := comp.keyHolds
	for k, new := range keyHolds {
		if new && !olds[k] {
			comp.anims[k].Reset()
		}
	}
	comp.keyHolds = keyHolds
}

func (comp HoldLightsComp) Draw(dst draws.Image) {
	for _, a := range comp.anims {
		a.Draw(dst)
	}
}
