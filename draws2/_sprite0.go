package draws

import "github.com/hajimehoshi/ebiten/v2"

type Position struct {
	Point
	Origin Origin
}

type Scale struct {
	Point
	Filter
}

const (
	FilterLinear = iota
	FilterNearest
)

type Filter int
type Sprite struct {
	i *ebiten.Image
	Position
	Scale
}

func NewSprite(path string) *Sprite {
	return NewSpriteFromImage(NewImage(path))
}
func NewSpriteFromImage(i *ebiten.Image) *Sprite {
	return &Sprite{
		i: i,
		Scale: Scale{
			Point:  Point{1, 1},
			Filter: FilterLinear,
		},
	}
}
func (s Sprite) Size() Point {
	if s.i == nil {
		return Point{}
	}
	return IntPt(s.i.Size()).Mul(s.Scale.Point)
}
func (s *Sprite) SetSize(size Point) {
	if s.i == nil {
		return
	}
	s.Scale.Point = size.Div(IntPt(s.i.Size()))
}
