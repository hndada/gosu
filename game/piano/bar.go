package piano

import (
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

// Drum and Piano modes have different bar drawing methods.
// Hence, this method is defined per game mode.
type Bar struct {
	position float64
}

type Bars struct {
	data  []Bar
	index int
}

// Given Dynamics' index is 0 and it is fine to modify it.
func NewBars(dys game.Dynamics) Bars {
	// const useDefaultMeter = 0
	times := dys.BeatTimes()
	bs := make([]Bar, len(times))
	dys.Reset()
	for i, t := range times {
		dys.UpdateIndex(t)
		bs[i] = Bar{position: dys.Position(t)}
	}
	return Bars{data: bs}
}

type BarsComponent struct {
	bars   Bars
	sprite draws.Sprite
	cursor float64
	reach  float64 // the distance from the top of the screen to the stage base position.
}

func NewBarsComponent(res *Resources, opts *Options, c *Chart) (cmp BarsComponent) {
	s := draws.NewSprite(res.BarImage)
	s.SetSize(opts.StageWidths[c.keyCount], opts.BarHeight)
	s.Locate(opts.StagePositionX, opts.KeyPositionY, draws.CenterBottom)
	cmp.sprite = s
	cmp.bars = NewBars(c.Dynamics)
	cmp.reach = opts.KeyPositionY
	return
}

// When speed changes from fast to slow, which means there are more bars
// on the screen, updateHighestBar() will handle it optimally.
// When speed changes from slow to fast, which means there are fewer bars
// on the screen, updateHighestBar() actually does nothing, which is still
// fine because that makes some unnecessary bars are drawn.
// The same concept also applies to notes.
func (cmp *BarsComponent) Update(cursor float64) {
	cmp.cursor = cursor
	lowermost := cursor - cmp.reach // game.ScreenH
	for i := cmp.bars.index; i < len(cmp.bars.data); i++ {
		b := cmp.bars.data[i]
		if b.position > lowermost {
			break
		}
		// index should be updated outside of if block.
		cmp.bars.index = i
	}
}

// Bars are fixed. Lane itself moves, all bars move as same amount.
func (cmp BarsComponent) Draw(dst draws.Image) {
	uppermost := cmp.cursor + cmp.reach // game.ScreenH
	for i := cmp.bars.index; i < len(cmp.bars.data); i++ {
		b := cmp.bars.data[i]
		if b.position > uppermost {
			break
		}

		s := cmp.sprite
		pos := b.position - cmp.cursor
		s.Move(0, -pos)
		s.Draw(dst)
	}
}
