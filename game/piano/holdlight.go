package piano

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type HoldLights struct {
	anims               []draws.Animation
	keysLongNoteHolding []bool
	notes               []Note
	keysFocusedNote     []int
}

func NewHoldLights(res *Resources, opts *Options, c *Chart) HoldLights {
	holdLights := HoldLights{}
	holdLights.anims = make([]draws.Animation, c.keyCount)
	xs := opts.keyPositionXsMap[c.keyCount]
	for k := range holdLights.anims {
		a := draws.NewAnimation(res.HoldLightsFrames, 300)
		a.Scale(opts.HoldLightImageScale)
		a.Locate(xs[k], opts.KeyPositionY-opts.HintHeight/2, draws.CenterMiddle)
		a.ColorScale.Scale(1, 1, 1, opts.HoldLightOpacity)
		holdLights.anims[k] = a
	}
	holdLights.keysLongNoteHolding = make([]bool, c.keyCount)
	holdLights.notes = c.notes
	holdLights.keysFocusedNote = c.keysFocusedNote
	return holdLights
}

// draws only when a long note is holding.
func (holdLights *HoldLights) Update(ka game.KeyboardAction) {
	kfns := make([]Note, len(holdLights.keysFocusedNote)) // key focused notes
	for k, ni := range holdLights.keysFocusedNote {
		if ni < 0 || ni == len(holdLights.notes) {
			continue
		}
		kfns[k] = holdLights.notes[ni]
	}

	keysOld := holdLights.keysLongNoteHolding
	keysNew := holdLights.newKeysLongNoteHolding(ka, kfns)
	for k, new := range keysNew {
		old := keysOld[k]
		if (old && !new) || (!old && new) {
			holdLights.anims[k].Reset()
		}
	}
	holdLights.keysLongNoteHolding = keysNew
}

func (holdLights HoldLights) newKeysLongNoteHolding(ka game.KeyboardAction, kn []Note) []bool {
	klnh := make([]bool, len(kn))
	for k, holding := range ka.KeysHolding() {
		if holding && kn[k].Kind == Tail {
			klnh[k] = true
		}
	}
	return klnh
}

func (holdLights HoldLights) Draw(dst draws.Image) {
	for k, a := range holdLights.anims {
		if holdLights.keysLongNoteHolding[k] {
			a.Draw(dst)
		}
	}
}
