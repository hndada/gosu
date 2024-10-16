package piano

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/plays"
)

type HoldLightsComponent struct {
	anims               []draws.Animation
	keysLongNoteHolding []bool
	notes               *Notes
}

func NewHoldLightsComponent(res *Resources, opts *Options, c *Chart) (cmp HoldLightsComponent) {
	cmp.anims = make([]draws.Animation, c.keyCount)
	xs := opts.keyPositionXsMap[c.keyCount]
	for k := range cmp.anims {
		a := draws.NewAnimation(res.HoldLightsFrames, 300)
		a.Scale(opts.HoldLightImageScale)
		a.Locate(xs[k], opts.KeyPositionY-opts.HintHeight/2, draws.CenterMiddle)
		a.ColorScale.Scale(1, 1, 1, opts.HoldLightOpacity)
		cmp.anims[k] = a
	}
	cmp.keysLongNoteHolding = make([]bool, c.keyCount)
	cmp.notes = &c.Notes
	return
}

// draws only when a long note is holding.
func (cmp *HoldLightsComponent) Update(ka plays.KeyboardAction) {
	kfns := make([]Note, cmp.notes.keyCount) // key focused notes
	for k, ni := range cmp.notes.keysFocus {
		if ni < 0 || ni == len(cmp.notes.data) {
			continue
		}
		kfns[k] = cmp.notes.data[ni]
	}

	keysOld := cmp.keysLongNoteHolding
	keysNew := cmp.newKeysLongNoteHolding(ka, kfns)
	for k, new := range keysNew {
		old := keysOld[k]
		if (old && !new) || (!old && new) {
			cmp.anims[k].Reset()
		}
	}
	cmp.keysLongNoteHolding = keysNew
}

func (cmp HoldLightsComponent) newKeysLongNoteHolding(ka plays.KeyboardAction, kn []Note) []bool {
	klnh := make([]bool, len(kn))
	for k, holding := range ka.KeysHolding() {
		if holding && kn[k].Kind == Tail {
			klnh[k] = true
		}
	}
	return klnh
}

func (cmp HoldLightsComponent) Draw(dst draws.Image) {
	for k, a := range cmp.anims {
		if cmp.keysLongNoteHolding[k] {
			a.Draw(dst)
		}
	}
}
