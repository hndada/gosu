package draws

import "github.com/hajimehoshi/ebiten/v2"

type sizer interface{ Size() Vector2 }

// All of these are for defining size and position.
type Rectangle struct {
	Base     *Rectangle
	Viewport *Rectangle
	src      sizer
	Size     Length2
	Theta    float64
	Position Length2
	Aligns   Aligns
}

func NewRectangle(src sizer) Rectangle {
	w, h := src.Size().XY()
	return Rectangle{
		Size: NewLength2(nil, w, h, Pixel),
	}
}

// func (r Rectangle) root() *Rectangle {
// 	if r.Parent == nil {
// 		return &r
// 	}
// 	return r.Parent.root()
// }

// type whxyKind int

// const (
// 	kindW whxyKind = iota
// 	kindH
// 	kindX
// 	kindY
// )

// func (r Rectangle) pixel(kind whxyKind) float64 {
// 	l := [4]Length{r.W, r.H, r.X, r.Y}[kind]
// 	switch l.Unit {
// 	case Pixel:
// 		return l.Value
// 	case Percent:
// 		if r.Parent != nil {
// 			ratio := l.Value / 100.0
// 			return ratio * r.Parent.pixel(kind)
// 		}
// 	case RootPercent:
// 		if root := r.root(); root != nil {
// 			ratio := l.Value / 100.0
// 			return ratio * root.pixel(kind)
// 		}
// 	case Extra:
// 		if r.Parent != nil {
// 			return r.Parent.pixel(kind) + l.Value
// 		}
// 	}
// 	return l.Value
// }

// func (r Rectangle) w() float64 { return r.pixel(kindW) }
// func (r Rectangle) h() float64 { return r.pixel(kindH) }
// func (r Rectangle) x() float64 { return r.pixel(kindX) }
// func (r Rectangle) y() float64 { return r.pixel(kindY) }

// func (r *Rectangle) addPixel(kind whxyKind, pixel float64) {
// 	l := [4]Length{r.W, r.H, r.X, r.Y}[kind]
// 	switch l.Unit {
// 	case Pixel:
// 		l.Value += pixel
// 	case Percent:
// 		if r.Parent != nil {
// 			ratio := pixel / r.Parent.pixel(kind)
// 			l.Value += ratio * 100.0
// 		}
// 	case RootPercent:
// 		if root := r.root(); root != nil {
// 			ratio := pixel / root.pixel(kind)
// 			l.Value += ratio * 100.0
// 		}
// 	case Extra:
// 		l.Value += pixel
// 	}
// }

// func (r *Rectangle) AddPixelToW(pixel float64) { r.addPixel(kindW, pixel) }
// func (r *Rectangle) AddPixelToH(pixel float64) { r.addPixel(kindH, pixel) }
// func (r *Rectangle) AddPixelToX(pixel float64) { r.addPixel(kindX, pixel) }
// func (r *Rectangle) AddPixelToY(pixel float64) { r.addPixel(kindY, pixel) }

// func (r Rectangle) Size() Vector2 { return Vec2(r.w(), r.h()) }
// func (r *Rectangle) SetSize(w, h float64, unit Unit) {
// 	r.W = Length{w, unit}
// 	r.H = Length{h, unit}
// }
func (r *Rectangle) Scale(x, y float64) { r.Size.Mul(Vec2(x, y)) }

// func (r Rectangle) Position() Vector2 { return Vec2(r.x(), r.y()) }
// func (r *Rectangle) SetPosition(x, y float64, unit Unit) {
// 	r.X = Length{x, unit}
// 	r.Y = Length{y, unit}
// }
func (r *Rectangle) Move(x, y float64) { r.Position.Add(Vec2(x, y)) }

// Min is the left-top position of the box.
func (r Rectangle) Min() Vector2 {
	min := r.Position.Pixel()
	w := r.Size.X.Pixel()
	h := r.Size.Y.Pixel()
	min.X -= []float64{0, w / 2, w}[r.Aligns.X]
	min.Y -= []float64{0, h / 2, h}[r.Aligns.Y]
	if r.Base != nil {
		min = min.Add(r.Base.Min())
		if vp := r.Base.Viewport; vp != nil {
			min = min.Sub(vp.Min())
		}
	}
	return min
}
func (r Rectangle) Max() Vector2 { return r.Min().Add(r.Size.Pixel()) }

func (r Rectangle) geoM() ebiten.GeoM {
	geom := ebiten.GeoM{}
	scale := r.Size.Pixel().Div(r.src.Size())
	geom.Scale(scale.XY())
	geom.Rotate(r.Theta)
	geom.Translate(r.Min().XY())
	return geom
}

func (r Rectangle) In(p Vector2) bool {
	min := r.Min()
	max := r.Max()
	return min.X <= p.X && p.X < max.X &&
		min.Y <= p.Y && p.Y < max.Y
}

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
