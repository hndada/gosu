package draws

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Text struct {
	Text   string
	Color  color.Color // color.NRGBA
	face   font.Face
	origin Origin // Text's align is determined by its Origin.
	align  int
	xy
	wh
}

func (t *Text) SetFace(f font.Face, origin Origin) {
	t.face = f
	t.origin = origin
	bound := text.BoundString(t.face, t.Text)
	t.w = float64(bound.Max.X)
	t.h = float64(-bound.Min.Y)
}

func (t Text) Draw(screen *ebiten.Image, b Box) {
	var x, y float64
	switch t.origin.PositionX() {
	case OriginLeft:
		x = b.Pad.W
	case OriginCenter:
		x = b.w/2 - t.w/2
	case OriginRight:
		x = b.w - t.w - b.Pad.W
	}
	switch t.origin.PositionY() {
	case OriginTop:
		y = b.Pad.H + t.h
	case OriginMiddle:
		y = b.h/2 + t.h/2
	case OriginBottom:
		y = b.h - b.Pad.H
	}
	x = math.Floor(b.LeftTopX() + x)
	y = math.Ceil(b.LeftTopY() + y)
	text.Draw(screen, t.Text, t.face, int(x), int(y), t.Color)
}
