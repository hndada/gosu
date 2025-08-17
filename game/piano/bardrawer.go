package piano

import (
	"github.com/hndada/gosu/draws"
)

type BarDrawer struct {
	bars   []Bar
	index  int
	sprite draws.Sprite
	cursor float64
	reach  float64 // the distance from the top of the screen to the stage base position.
}

func NewBarDrawer(res *Resources, opts *Options, c *Chart) BarDrawer {
	bd := BarDrawer{}
	bd.bars = c.bars

	s := draws.NewSprite(res.BarImage)
	s.SetSize(opts.StageWidths[c.keyCount], opts.BarHeight)
	s.Locate(opts.StagePositionX, opts.KeyPositionY, draws.CenterBottom)
	bd.sprite = s

	bd.reach = opts.KeyPositionY
	return bd
}

// When speed changes from fast to slow, which means there are more bars
// on the screen, updateHighestBar() will handle it optimally.
// When speed changes from slow to fast, which means there are fewer bars
// on the screen, updateHighestBar() actually does nothing, which is still
// fine because that makes some unnecessary bars are drawn.
// The same concept also applies to notes.
func (bd *BarDrawer) Update(cursor float64) {
	bd.cursor = cursor
	lowermost := cursor - bd.reach // game.ScreenH
	for i := bd.index; i < len(bd.bars); i++ {
		b := bd.bars[i]
		if b.position > lowermost {
			break
		}
		// index should be updated outside of if block.
		bd.index = i
	}
}

// Bars are fixed. Lane itself moves, all bars move as same amount.
func (bd BarDrawer) Draw(dst draws.Image) {
	uppermost := bd.cursor + bd.reach // game.ScreenH
	for i := bd.index; i < len(bd.bars); i++ {
		b := bd.bars[i]
		if b.position > uppermost {
			break
		}

		s := bd.sprite
		pos := b.position - bd.cursor
		s.Move(0, -pos)
		s.Draw(dst)
	}
}
