package piano

import (
	"image/color"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

// Config is for wrapping required arguments.
type BarConfig struct {
	positionX float64 // from FieldPositionX
	positionY float64 // from KeyPositionY
	width     float64 // from FieldWidth
	Height    float64
}

func NewBarConfig(screen mode.ScreenConfig, stage StageConfig) BarConfig {
	return BarConfig{
		positionX: stage.FieldPositionX,
		positionY: stage.HintPositionY,
		width:     stage.Width(),
		Height:    1,
	}
}

// bottom: hit position
type BarComponent struct {
	bars    []*mode.Bar
	sprite  draws.Sprite
	highest *mode.Bar
	cursor  float64
}

func NewBarComponent(cfg BarConfig, bars []*mode.Bar) (bc BarComponent) {
	bc.bars = bars

	// Bar component uses a simple white rectangle as sprite.
	img := draws.NewImage(cfg.width, cfg.Height)
	img.Fill(color.White)
	sprite := draws.NewSprite(img)
	sprite.Locate(cfg.positionX, cfg.positionY, draws.CenterMiddle)
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
