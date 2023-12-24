package piano

import (
	"image/color"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

// Bar component uses a simple white rectangle as sprite.
type BarRes struct {
	// Bar component requires no external resources.
}

type BarOpts struct {
	w float64
	H float64
	x float64
	y float64 // center bottom
}

func NewBarOpts(keys KeysOpts) BarOpts {
	return BarOpts{
		w: keys.stageW,
		H: 1,
		x: keys.StageX,
		y: keys.BaselineY,
	}
}

type BarComp struct {
	bars   []mode.Bar
	idx    int
	cursor float64
	sprite draws.Sprite
}

func NewBarComp(res BarRes, opts BarOpts, bars []mode.Bar) (comp BarComp) {
	comp.bars = bars

	img := draws.NewImage(opts.w, opts.H)
	img.Fill(color.White)

	sprite := draws.NewSprite(img)
	sprite.Locate(opts.x, opts.y, draws.CenterBottom)
	comp.sprite = sprite
	return
}

// When speed changes from fast to slow, which means there are more bars
// on the screen, updateHighestBar() will handle it optimally.
// When speed changes from slow to fast, which means there are fewer bars
// on the screen, updateHighestBar() actually does nothing, which is still
// fine because that makes some unnecessary bars are drawn.
// The same concept also applies to notes.
func (comp *BarComp) Update(cursor float64) {
	lowerBound := cursor - mode.ScreenH
	for i := comp.idx; i < len(comp.bars); i++ {
		b := comp.bars[i]
		if b.Position > lowerBound {
			break
		}
		// index should be updated outside of if block.
		comp.idx = i
	}
	comp.cursor = cursor
}

// Bars are fixed. Lane itself moves, all bars move as same amount.
func (comp BarComp) Draw(dst draws.Image) {
	upperBound := comp.cursor + mode.ScreenH
	for i := comp.idx; i < len(comp.bars); i++ {
		b := comp.bars[i]
		if b.Position > upperBound {
			break
		}

		pos := b.Position - comp.cursor
		sprite := comp.sprite
		sprite.Move(0, -pos)
		sprite.Draw(dst)
	}
}
