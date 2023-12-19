package piano

import (
	"image/color"

	"github.com/hndada/gosu/draws"
	mode "github.com/hndada/gosu/mode2"
)

// Since XxxConfig is for embedding, it's better
// not to omit Xxx prefix on each field.
type BarConfig struct {
	BarHeight float64
}

func NewBarConfig() BarConfig {
	return BarConfig{
		BarHeight: 1,
	}
}

// Arguments is used only for NewBarComponent,
// hence omitting Xxx prefix is acceptable.
type BarArgs struct {
	PositionX float64
	PositionY float64
	Width     float64
	Height    float64
}

func (cfg Config) BarArgs() BarArgs {
	return BarArgs{
		PositionX: cfg.StagePositionX,
		PositionY: cfg.HitPositionY,
		Width:     cfg.StageWidth(),
		Height:    cfg.BarHeight,
	}
}

// bottom: hit position
type BarComponent struct {
	bars    []*mode.Bar
	sprite  draws.Sprite
	highest *mode.Bar
	cursor  float64
}

func NewBarComponent(bars []*mode.Bar, args BarArgs) (bc BarComponent) {
	bc.bars = bars

	// Bar component uses a simple white rectangle as sprite.
	img := draws.NewImage(args.Width, args.Height)
	img.Fill(color.White)

	sprite := draws.NewSprite(img)
	sprite.Locate(args.PositionX, args.PositionY, draws.CenterMiddle)
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
