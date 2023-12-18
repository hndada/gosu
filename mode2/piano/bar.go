package piano

import (
	"image/color"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

// bottom: hit position
type BarComponent struct {
	bars    mode.Bars
	sprite  draws.Sprite
	highest *mode.Bar
	cursor  float64
}

func NewBarComponent(bars mode.Bars, cfg BarConfig) (bc BarComponent) {
	bc.bars = bars

	img := draws.NewImage(*cfg.fieldWidth, cfg.Height)
	img.Fill(color.White)
	sprite := draws.NewSprite(img)
	sprite.Locate(*cfg.fieldPositionX, *cfg.keyPositionY, draws.CenterBottom)
	bc.sprite = sprite
	return
}
func (bc *BarComponent) Update(cursor float64) {
	bc.highest = bc.bars.Highest(cursor)
	bc.cursor = cursor

}

// Bars are fixed. Lane itself moves, all bars move as same amount.
func (bc BarComponent) Draw(dst draws.Image) {
	lowerBound := bc.cursor - 100
	for b := bc.highest; b != nil && b.Position > lowerBound; b = b.Prev {
		pos := b.Position - bc.cursor
		sprite := bc.sprite
		sprite.Move(0, -pos)
		sprite.Draw(dst, draws.Op{})
		if b.Prev == nil {
			break
		}
	}
}
