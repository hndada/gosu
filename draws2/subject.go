package draws

import "github.com/hajimehoshi/ebiten/v2"

type Subject interface {
	Size() Point
	// SetSize(w, h float64)
	Draw(*ebiten.Image, ebiten.DrawImageOptions)
	In(p Point) bool
}
