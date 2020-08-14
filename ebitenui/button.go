package ebitenui

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

type Button struct {
	MinPt     image.Point
	Image     *ebiten.Image
	mouseDown bool
	onPressed func(b *Button) // memo: an example that field of a struct controls the whole struct

	// Padding   image.Point // might be better that input image is already padded
}

func (b *Button) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p := image.Pt(ebiten.CursorPosition())
		if p.In(b.rect()) {
			b.mouseDown = true
		} else {
			b.mouseDown = false
		}
	} else { // onPressed should not be called when user is still pressing
		if b.mouseDown {
			if b.onPressed != nil {
				b.onPressed(b)
			}
		}
		b.mouseDown = false
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.MinPt.X), float64(b.MinPt.Y))
	screen.DrawImage(b.Image, op)
}

func (b *Button) SetOnPressed(f func(b *Button)) {
	b.onPressed = f
}

func (b *Button) rect() image.Rectangle {
	w, h := b.Image.Size()
	maxPt := b.MinPt.Add(image.Pt(w, h))
	return image.Rectangle{b.MinPt, maxPt}
}