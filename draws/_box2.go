package draws

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// Sprite is scaled image.
// Default filter is Linear. Use Op's Filter want other.
type Sprite2 struct {
	Position
	size Size
	i    *ebiten.Image
}
type Label struct {
	Position
	size  Size
	Text  string
	Face  font.Face
	Color color.Color
}

// Box is for wrapping sprites and labels.
type Box struct {
}

// type Container interface {}

// A box may have one text and one image.
// type Box struct {
// 	// Prev, Next                    *Box
// 	// Parent, FirstChild, LastChild *Box
// 	Sprite
// 	Pad WH
// 	// Margin WH
// 	Text
// }

// https://www.w3schools.com/css/css_grid.asp
// type Box2 struct {
// 	XY
// 	Pad WH
// 	Origin
// }

// If Sprite or Text is a item of Grid, their x and y are ignored.
type Grid struct {
	Rows, Cols int
}
type GridItem struct {
	Row, Col int
	// Span     struct{ W, H int }
	Item interface {
		Draw(screen *ebiten.Image)
	}
}
