package draws

import "github.com/hajimehoshi/ebiten/v2"

// All of these are for defining size and position.
type Rectangle struct {
	Parent   *Rectangle
	Viewport *Rectangle
	W, H     Length
	Theta    float64
	X, Y     Length
	Aligns   Aligns
}

func NewRectangle(w, h float64) Rectangle {
	return Rectangle{
		W: Length{w, Pixel},
		H: Length{h, Pixel},
	}
}

type whxyKind int

const (
	kindW whxyKind = iota
	kindH
	kindX
	kindY
)

func (r Rectangle) pixel(kind whxyKind) float64 {
	l := [4]Length{r.W, r.H, r.X, r.Y}[kind]
	switch l.Unit {
	case Pixel:
		return l.Value
	case Percent:
		if r.Parent != nil {
			return l.Value * r.Parent.pixel(kind) / 100.0
		}
	case RootPercent:
		if root := r.root(); root != nil {
			return l.Value * root.pixel(kind) / 100.0
		}
	case Extra:
		if r.Parent != nil {
			return r.Parent.pixel(kind) + l.Value
		}
	}
	return l.Value
}

func (r Rectangle) root() *Rectangle {
	if r.Parent == nil {
		return &r
	}
	return r.Parent.root()
}

func (r Rectangle) w() float64 { return r.pixel(kindW) }
func (r Rectangle) h() float64 { return r.pixel(kindH) }
func (r Rectangle) x() float64 { return r.pixel(kindX) }
func (r Rectangle) y() float64 { return r.pixel(kindY) }

func (r Rectangle) Size() Vector2 { return Vec2(r.w(), r.h()) }
func (r *Rectangle) SetSize(w, h float64, unit Unit) {
	r.W = Length{w, unit}
	r.H = Length{h, unit}
}
func (r *Rectangle) Scale(scale float64) {
	r.W.Value *= scale
	r.H.Value *= scale
}

func (r Rectangle) Position() Vector2 { return Vec2(r.x(), r.y()) }
func (r *Rectangle) SetPosition(x, y float64, unit Unit) {
	r.X = Length{x, unit}
	r.Y = Length{y, unit}
}
func (r *Rectangle) Move(x, y float64) {
	r.X.Value += x
	r.Y.Value += y
}

func (r Rectangle) geoM(srcSize Vector2) ebiten.GeoM {
	geom := ebiten.GeoM{}
	scale := r.Size().Div(srcSize)
	geom.Scale(scale.XY())
	geom.Rotate(r.Theta)
	geom.Translate(r.Min().XY())
	return geom
}

// Min is the left-top position of the box.
func (r Rectangle) Min() Vector2 {
	min := r.Position()
	min.X -= []float64{0, r.w() / 2, r.w()}[r.Aligns.X]
	min.Y -= []float64{0, r.h() / 2, r.h()}[r.Aligns.Y]
	if r.Parent != nil {
		min = min.Add(r.Parent.Min())
	}
	if r.Viewport != nil {
		min = min.Sub(r.Viewport.Min())
	}
	return min
}
func (r Rectangle) Max() Vector2 { return r.Min().Add(r.Size()) }

func (r Rectangle) Exposed(dst Image) bool {
	return dst.In(r.Min()) || dst.In(r.Max())
}

// func (r Rectangle) In(p Vector2) bool {
// 	min := r.Min()
// 	max := r.Max()
// 	return min.X <= p.X && p.X < max.X &&
// 		min.Y <= p.Y && p.Y < max.Y
// }

// func (r1 Rectangle) Intersect(r2 Rectangle) bool {
// 	return r1.In(r2.Min()) || r1.In(r2.Max())
// }

// Rectangle implements Shape.
// type Shape interface {
// 	Size() Vector2
// 	Position() Vector2
// 	Min() Vector2
// 	Max() Vector2
// }
