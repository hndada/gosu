package draws

import "github.com/hajimehoshi/ebiten/v2"

type Subject interface {
	Size() Point
	SetSize(size Point)
	Draw(*ebiten.Image, ebiten.DrawImageOptions)
	In(Point) bool
}
