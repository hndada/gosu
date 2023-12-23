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
	stage    draws.WHXY // parent element
	baseline draws.Position
	H        float64
}

func NewBarOpts(stage draws.WHXY, baseline draws.Position) BarOpts {
	return BarOpts{
		stage:    stage,
		baseline: baseline,
		H:        1,
	}
}

type BarComp struct {
	bars    []*mode.Bar
	sprite  draws.Sprite
	highest *mode.Bar
	cursor  float64
}

func NewBarComp(res BarRes, opts BarOpts, bars []*mode.Bar) (comp BarComp) {
	comp.bars = bars

	img := draws.NewImage(opts.stage.X, opts.H)
	img.Fill(color.White)

	sprite := draws.NewSprite(img)
	sprite.Locate(opts.stage.X, opts.baseline.Y, draws.CenterMiddle)
	comp.sprite = sprite
	return
}

func (comp *BarComp) Update(cursor float64) {
	comp.highest = comp.bars.Highest(cursor)
	comp.cursor = cursor
}

// Bars are fixed. Lane itself moves, all bars move as same amount.
func (comp BarComp) Draw(dst draws.Image) {
	lowerBound := comp.cursor - 100
	for b := comp.highest; b != nil && b.Position > lowerBound; b = b.Prev {
		pos := b.Position - comp.cursor
		sprite := comp.sprite
		sprite.Move(0, -pos)
		sprite.Draw(dst, draws.Op{})
		if b.Prev == nil {
			break
		}
	}
}
