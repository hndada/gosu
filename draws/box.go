package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Box struct {
	Point
	Origin
	Subject
}

func (b Box) Min() Point {
	p := b.Point
	w, h := b.WH().XY()
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
	w, h := b.WH().XY()
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
	w, h := b.WH().XY()
	return p.X >= 0 && p.X <= w && p.Y >= 0 && p.Y <= h
}
func (b Box) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	min := b.Min()
	b.Subject.Draw(screen, min.X, min.Y, op)
}

// func (b Box) minX() float64 {
// 	switch b.Origin.PositionX() {
// 	case OriginLeft:
// 		return b.X
// 	case OriginCenter:
// 		return b.X - b.W/2
// 	case OriginRight:
// 		return b.X - b.W
// 	}
// 	return 0
// }
// func (b Box) minY() float64 {
// 	switch b.Origin.PositionY() {
// 	case OriginTop:
// 		return b.Y
// 	case OriginMiddle:
// 		return b.Y - b.H/2
// 	case OriginBottom:
// 		return b.Y - b.H
// 	}
// 	return 0
// }

// func (b Box) Max() Point {
// 	return Point{b.maxX(), b.maxY()}
// }

// func (b Box) maxX() float64 {
// 	switch b.Origin.PositionX() {
// 	case OriginLeft:
// 		return b.X + b.W
// 	case OriginCenter:
// 		return b.X + b.W/2
// 	case OriginRight:
// 		return b.X
// 	}
// 	return 0
// }
// func (b Box) maxY() float64 {
// 	switch b.Origin.PositionY() {
// 	case OriginTop:
// 		return b.Y + b.H
// 	case OriginMiddle:
// 		return b.Y + b.H/2
// 	case OriginBottom:
// 		return b.Y
// 	}
// 	return 0
// }
