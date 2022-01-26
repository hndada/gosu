package ui

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type LongSprite struct {
	Sprite
	Vertical bool
}

// temp: no need to be method of LongSprite, to make sure only LongSprite uses this
func (s LongSprite) isOut(w, h, x, y int, screenSize image.Point) bool {
	return x+w < 0 || x > screenSize.X || y+h < 0 || y > screenSize.Y
}

// A long image should be drawn in pieces; there's height limit in *ebiten.Image
func (s LongSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w1, h1 := s.i.Size()
	switch s.Vertical {
	case true:
		op.GeoM.Scale(s.scaleW(), 1) // height 쪽은 굳이 scale 하지 않는다
		// important: op is not AB = BA
		x, y := s.X, s.Y
		op.GeoM.Translate(float64(x), float64(y))
		q, r := s.H/h1, s.H%h1+1 // quotient, remainder // temp: +1

		first := s.i.Bounds()
		w, h := s.W, r
		first.Min = image.Pt(0, h1-r)
		if !s.isOut(w, h, x, y, screen.Bounds().Size()) {
			screen.DrawImage(s.i.SubImage(first).(*ebiten.Image), op)
		}
		op.GeoM.Translate(0, float64(h))
		y += h
		h = h1
		for i := 0; i < q; i++ {
			if !s.isOut(w, h, x, y, screen.Bounds().Size()) {
				screen.DrawImage(s.i, op)
			}
			op.GeoM.Translate(0, float64(h))
			y += h
		}

	default:
		op.GeoM.Scale(1, s.scaleH())
		op.GeoM.Translate(float64(s.X), float64(s.Y))
		q, r := s.W/w1, s.W%w1+1 // temp: +1

		for i := 0; i < q; i++ {
			screen.DrawImage(s.i, op)
			op.GeoM.Translate(float64(w1), 0)
		}

		last := s.i.Bounds()
		last.Max = image.Pt(r, h1)
		screen.DrawImage(s.i.SubImage(last).(*ebiten.Image), op)
	}
}
