package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Box is wrapped Subject(Sprite, Label) with Position data.
// x and y of Sub box are determined by Outer box.
type Box struct {
	Subject
	Point
	Origin
	Align
	Pad  Point
	Gap  Point // Margin
	Subs [][]Box
}

func (b Box) Min() Point {
	p := b.Point
	w, h := b.Size().XY()
	switch b.Origin.PositionX() {
	case OriginLeft:
		p.X -= 0
	case OriginCenter:
		p.X -= w / 2
	case OriginRight:
		p.X -= w
	}
	switch b.Origin.PositionY() {
	case OriginTop:
		p.Y -= 0
	case OriginMiddle:
		p.Y -= h / 2
	case OriginBottom:
		p.Y -= h
	}
	return p
}
func (b Box) Max() Point {
	p := b.Point
	w, h := b.Size().XY()
	switch b.Origin.PositionX() {
	case OriginLeft:
		p.X += w
	case OriginCenter:
		p.X += w / 2
	case OriginRight:
		p.X += 0
	}
	switch b.Origin.PositionY() {
	case OriginTop:
		p.Y += h
	case OriginMiddle:
		p.Y += h / 2
	case OriginBottom:
		p.Y += 0
	}
	return p
}

// Input is usually cursor's position.
func (b Box) In(p Point) bool {
	p = p.Sub(b.Min())
	w, h := b.Size().XY()
	return p.X >= 0 && p.X <= w && p.Y >= 0 && p.Y <= h
}
func (b Box) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	b.Subject.Draw(screen, op, b.Min())
}

// https://www.w3schools.com/css/css_grid.asps
// Star model: Prev, Next, Parent, First/Last Child
