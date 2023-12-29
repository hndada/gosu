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

func NewBars(dys game.Dynamics) []Bar {
	// const useDefaultMeter = 0
	times := dys.BeatTimes()
	bs := make([]Bar, 0, len(times))
	for _, t := range times {
		b := Bar{Time: t}
		bs = append(bs, b)
	}

	// Set bar positions
	// dys.Index = 0
	for i, b := range bs {
		dys.UpdateIndex(b.Time)
		d := dys.Current()
		bs[i].position = d.Position + float64(b.Time-d.Time)*d.Speed
	}
	return bs
}

// Drum and Piano modes have different bar drawing methods.
// Hence, this method is not defined in mode.go.

type BarRes struct {
	img draws.Image
}

func (res *BarRes) Load(fsys fs.FS) {
	// Uses generated image.
	img := draws.NewImage(1, 1)
	img.Fill(color.White)
	res.img = img
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
	bars   []Bar
	cursor float64
	lowest int
	sprite draws.Sprite
}

func NewBarComp(res BarRes, opts BarOpts, bars []Bar) (comp BarComp) {
	comp.bars = bars

	s := draws.NewSprite(res.img)
	s.SetSize(opts.w, opts.H)
	s.Locate(opts.x, opts.y, draws.CenterBottom)
	comp.sprite = s
	return
}

// When speed changes from fast to slow, which means there are more bars
// on the screen, updateHighestBar() will handle it optimally.
// When speed changes from slow to fast, which means there are fewer bars
// on the screen, updateHighestBar() actually does nothing, which is still
// fine because that makes some unnecessary bars are drawn.
// The same concept also applies to notes.
func (comp *BarComp) Update(cursor float64) {
	comp.cursor = cursor
	lowerBound := cursor - game.ScreenH
	for i := comp.lowest; i < len(comp.bars); i++ {
		b := comp.bars[i]
		if b.position > lowerBound {
			break
		}
		// index should be updated outside of if block.
		comp.lowest = i
	}
}

// Bars are fixed. Lane itself moves, all bars move as same amount.
func (comp BarComp) Draw(dst draws.Image) {
	upperBound := comp.cursor + game.ScreenH
	for i := comp.lowest; i < len(comp.bars); i++ {
		b := comp.bars[i]
		if b.position > upperBound {
			break
		}

		s := comp.sprite
		pos := b.position - comp.cursor
		s.Move(0, -pos)
		s.Draw(dst)
	}
}
