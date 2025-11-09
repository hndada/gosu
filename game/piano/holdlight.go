package piano

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

type HoldLights struct {
	anims            []draws.Animation
	longNoteHoldings []bool
	notes            []Note
	focusedNotes     []int
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
	holdLights.longNoteHoldings = make([]bool, c.keyCount)
	holdLights.notes = c.notes
	holdLights.focusedNotes = c.focusedNotes
	return holdLights
}

// draws only when a long note is holding.
func (holdLights *HoldLights) Update(ka game.KeyboardAction) {
	kfns := make([]Note, len(holdLights.focusedNotes)) // key focused notes
	for k, ni := range holdLights.focusedNotes {
		if ni < 0 || ni == len(holdLights.notes) {
			continue
		}
		kfns[k] = holdLights.notes[ni]
	}

	olds := holdLights.longNoteHoldings
	news := holdLights.newlongNoteHoldings(ka, kfns)
	for k, new := range news {
		old := olds[k]
		if (old && !new) || (!old && new) {
			holdLights.anims[k].Reset()
		}
	}
	holdLights.longNoteHoldings = news
}

func (holdLights HoldLights) newlongNoteHoldings(ka game.KeyboardAction, kn []Note) []bool {
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
		if holdLights.longNoteHoldings[k] {
			a.Draw(dst)
		}
	}
}
