package draws

import "github.com/hajimehoshi/ebiten/v2"

type Subject interface {
	Size() Vector2
	Draw(*ebiten.Image, ebiten.DrawImageOptions)
	In(Vector2) bool
}
