package ui

import "github.com/hajimehoshi/ebiten/v2"

type FixedSprite struct { // a sprite that never moves once appears
	sprite Sprite
	op     *ebiten.DrawImageOptions
}

func NewFixedSprite(s Sprite) FixedSprite {
	return FixedSprite{
		sprite: s,
		op:     s.Op(),
	}
}
func (s FixedSprite) Draw(screen *ebiten.Image) {
	if s.sprite.i == nil {
		panic("no image")
	}
	if s.sprite.isOut(screen.Bounds().Max) {
		return
	}
	screen.DrawImage(s.sprite.i, s.op)
}
