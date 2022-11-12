package draws

import "github.com/hajimehoshi/ebiten/v2"

type Subject interface {
	// Size() Vector2
	In(Vector2) bool
	Draw(*ebiten.Image, ebiten.DrawImageOptions)
}
