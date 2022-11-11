package draws

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite is an image drawn in a screen based on its position and scale.
// DrawImageOptions is not commutative. Do Translate at the final stage.
type Sprite struct {
	i *ebiten.Image
	Rectangle
}

func NewSprite(path string) Sprite {
	i := NewImage(path)
	return NewSpriteFromImage(i)
}
func NewSpriteFromImage(i *ebiten.Image) Sprite {
	if i == nil {
		return Sprite{}
	}
	return Sprite{
		i:         i,
		Rectangle: NewRectangle(ImageSize(i)),
	}
}
func (s Sprite) IsValid() bool { return s.i != nil }
func (s Sprite) Draw(screen *ebiten.Image, op ebiten.DrawImageOptions) {
	if s.i == nil {
		return
	}
	op.GeoM.Scale(s.Scale.XY())
	leftTop := s.LeftTop(ImageSize(screen))
	op.GeoM.Translate(leftTop.XY())
	op.Filter = s.Filter
	screen.DrawImage(s.i, &op)
}
