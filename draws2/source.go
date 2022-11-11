package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Source interface {
	IsValid() bool
	Size() Vector2
	Draw(screen *ebiten.Image, op ebiten.DrawImageOptions)
}
