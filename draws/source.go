package draws

import "github.com/hajimehoshi/ebiten/v2"

type Op = ebiten.DrawImageOptions

// Image and Text implements Source.
type Source interface {
	IsValid() bool
	Size() Vector2
	Draw(Image, Op)
}
