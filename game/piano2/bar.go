package piano

import (
	draws "github.com/hndada/gosu/draws5"
	"github.com/hndada/gosu/game"
)

// Drum and Piano modes have different bar drawing methods.
// Hence, this method is defined per game mode.
type Bar struct {
	position float64
}

type Bars struct {
	bars  []Bar
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
	return Bars{bars: bs}
}

type BarsOptions struct {
	w float64
	H float64
	x float64
	y float64
}

func NewBarsOptions(stage StageOptions) BarsOptions {
	return BarsOptions{
		w: stage.w,
		H: 1,
		x: stage.X,
		y: stage.H,
	}
}

type BarsComponent struct {
	sprite draws.Sprite
	Bars
	cursor float64
}

func NewBarsComponent(res BarsResources, opts BarsOptions, dys game.Dynamics) (cmp BarsComponent) {
	s := draws.NewSprite(res.img)
	s.SetSize(opts.w, opts.H)
	s.Locate(opts.x, opts.y, draws.CenterBottom)
	cmp.sprite = s
	cmp.Bars = NewBars(dys)
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
	lowermost := cursor - game.ScreenH
	for i := cmp.index; i < len(cmp.bars); i++ {
		b := cmp.bars[i]
		if b.position > lowermost {
			break
		}
		// index should be updated outside of if block.
		cmp.index = i
	}
}

// Bars are fixed. Lane itself moves, all bars move as same amount.
func (cmp BarsComponent) Draw(dst draws.Image) {
	uppermost := cmp.cursor + game.ScreenH
	for i := cmp.index; i < len(cmp.bars); i++ {
		b := cmp.bars[i]
		if b.position > uppermost {
			break
		}

		s := cmp.sprite
		pos := b.position - cmp.cursor
		s.Move(0, -pos)
		s.Draw(dst)
	}
}
