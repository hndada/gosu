package piano

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type HoldLightsResources struct {
	frames draws.Frames
}

func (br *HoldLightsResources) Load(fsys fs.FS) {
	fname := "piano/light/hold.png"
	br.frames = draws.NewFramesFromFile(fsys, fname)
}

type HoldLightsOptions struct {
	Scale    float64
	keyCount int
	keysX    []float64
	y        float64
	Opacity  float32
}

func NewHoldLightsOptions(keys KeysOptions) HoldLightsOptions {
	return HoldLightsOptions{
		Scale:    1.0,
		keyCount: keys.keyCount,
		keysX:    keys.x,
		y:        keys.y,
		Opacity:  1.2,
	}
}

type HoldLightsComponent struct {
	anims               []draws.Animation
	keysLongNoteHolding []bool
}

func NewHoldLightsComponent(res HoldLightsResources, opts HoldLightsOptions) (cmp HoldLightsComponent) {
	cmp.anims = make([]draws.Animation, opts.keyCount)
	for k := range cmp.anims {
		a := draws.NewAnimation(res.frames, 300)
		a.MultiplyScale(opts.Scale)
		a.Locate(opts.keysX[k], opts.y, draws.CenterBottom)
		a.ColorScale.Scale(1, 1, 1, opts.Opacity)
		cmp.anims[k] = a
	}
	return
}

func (cmp *HoldLightsComponent) Update(ka game.KeyboardAction, kn []Note) {
	keysOld := cmp.keysLongNoteHolding
	keysNew := cmp.newKeysLongNoteHolding(ka, kn)
	for k, new := range keysNew {
		old := keysOld[k]
		if !old && new {
			cmp.anims[k].Reset()
		}
	}
	cmp.keysLongNoteHolding = keysNew
}

func (cmp HoldLightsComponent) newKeysLongNoteHolding(ka game.KeyboardAction, kn []Note) []bool {
	klnh := make([]bool, len(kn))
	for k, holding := range ka.KeysHolding() {
		if holding && kn[k].Kind == Tail {
			klnh[k] = true
		}
	}
	return klnh
}

func (cmp HoldLightsComponent) Draw(dst draws.Image) {
	for _, a := range cmp.anims {
		a.Draw(dst)
	}
}
