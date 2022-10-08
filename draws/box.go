package draws

import "github.com/hajimehoshi/ebiten/v2"

// A box may have one text and one image.
type Box struct {
	// Prev, Next                    *Box
	// Parent, FirstChild, LastChild *Box
	Sprite
	Pad WH
	// Margin WH
	Text
}

func (b *Box) CursorIn() bool {
	cx, cy := ebiten.CursorPosition()
	return b.Sprite.In(float64(cx), float64(cy))
}

// https://www.w3schools.com/css/css_grid.asp
