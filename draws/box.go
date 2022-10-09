package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Box model
// Box is wrapped Subject(Sprite, Label) with Position data.
type Box struct {
	Subject
	Point
	Origin

	Subs    [][]Box
	Padding Point
	Margin  Point
	Align
}

// Zero-valued Point infers tight box.
func (b Box) mins() [][]Point {

}

// x and y of Sub box are determined by Outer box.
// Suppose boxes in a same row has equal height.
// Wait, no, labels have different height.
func (b Box) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	b.Subject.Draw(screen, op, b.Min())
	var offset Point
	for _, row := range b.Subs {
		offset.X = 0
		offset.Y += b.Pad.Y
		for _, sub := range row {
			offset.X += b.Pad.X
			sub.Point.Add(offset)
			sub.Draw(screen, op)
			offset.X += sub.Subject.Size().X
		}
		offset.Y += row[0].Subject.Size().Y
	}
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

// Input is usually cursor's position.
func (b Box) In(p Point) bool {
	p = p.Sub(b.Min())
	w, h := b.Size().XY()
	return p.X >= 0 && p.X <= w && p.Y >= 0 && p.Y <= h
}

// func (b Box) Max() Point {
// 	p := b.Point
// 	w, h := b.Size().XY()
// 	switch b.Origin.PositionX() {
// 	case OriginLeft:
// 		p.X += w
// 	case OriginCenter:
// 		p.X += w / 2
// 	case OriginRight:
// 		p.X += 0
// 	}
// 	switch b.Origin.PositionY() {
// 	case OriginTop:
// 		p.Y += h
// 	case OriginMiddle:
// 		p.Y += h / 2
// 	case OriginBottom:
// 		p.Y += 0
// 	}
// 	return p
// }
