package piano

import (
	"image/color"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
)

// Defining Bars is just redundant if it has no additional methods.
type Bar struct {
	Time     int32 // in milliseconds
	position float64
}

// Given Dynamics' index is 0 and it is fine to modify it.
func NewBars(dys game.Dynamics) []Bar {
	// const useDefaultMeter = 0
	times := dys.BeatTimes()
	bs := make([]Bar, 0, len(times))
	for _, t := range times {
		b := Bar{Time: t}
		bs = append(bs, b)
	}

	for i, b := range bs {
		dys.UpdateIndex(b.Time)
		d := dys.Current()
		bs[i].position = d.Position + float64(b.Time-d.Time)*d.Speed
	}
	return bs
}

// Drum and Piano modes have different bar drawing methods.
// Hence, this method is not defined in mode.go.
type BarsResources struct {
	img draws.Image
}

func (res *BarsResources) Load(fsys fs.FS) {
	// Uses generated image.
	img := draws.NewImage(1, 1)
	img.Fill(color.White)
	res.img = img
}

type BarsOptions struct {
	w float64
	H float64
	x float64
	y float64 // center bottom
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
	bars   []Bar
	cursor float64
	lowest int
	sprite draws.Sprite
}

func NewBarComponent(res BarsResources, opts BarsOptions, bars []Bar) (cmp BarsComponent) {
	cmp.bars = bars

	s := draws.NewSprite(res.img)
	s.SetSize(opts.w, opts.H)
	s.Locate(opts.x, opts.y, draws.CenterBottom)
	cmp.sprite = s
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
	lowerBound := cursor - game.ScreenH
	for i := cmp.lowest; i < len(cmp.bars); i++ {
		b := cmp.bars[i]
		if b.position > lowerBound {
			break
		}
		// index should be updated outside of if block.
		cmp.lowest = i
	}
}

// Bars are fixed. Lane itself moves, all bars move as same amount.
func (cmp BarsComponent) Draw(dst draws.Image) {
	upperBound := cmp.cursor + game.ScreenH
	for i := cmp.lowest; i < len(cmp.bars); i++ {
		b := cmp.bars[i]
		if b.position > upperBound {
			break
		}

		s := cmp.sprite
		pos := b.position - cmp.cursor
		s.Move(0, -pos)
		s.Draw(dst)
	}
}
