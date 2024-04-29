package piano

import (
	draws "github.com/hndada/gosu/draws6"
	"github.com/hndada/gosu/game"
)

type HoldLightsComponent struct {
	anims               []draws.Animation
	keysLongNoteHolding []bool
}

func NewHoldLightsComponent(res *Resources, opts *Options, keyCount int) (cmp HoldLightsComponent) {
	cmp.anims = make([]draws.Animation, keyCount)
	xs := opts.keyPositionXsMap[keyCount]
	for k := range cmp.anims {
		a := draws.NewAnimation(res.HoldLightsFrames, 300)
		a.Scale(opts.HoldLightImageScale)
		a.Locate(xs[k], opts.KeyPositionY, draws.CenterBottom)
		a.ColorScale.Scale(1, 1, 1, opts.HoldLightOpacity)
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
