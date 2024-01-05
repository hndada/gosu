package draws

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sprite, Label, Animation, Filler implement Drawable.
type Drawable interface {
	Draw(dst Image, imgOp *ebiten.DrawImageOptions)
}

// Separate types are required to use Source's methods.
type Sprite Box[Image]

func NewSprite(img Image) Sprite {

}

// l := Box[Text]{}
// l.Source.Text = "Hello, World!" // It works!
type Label Box[Text]

func NewLabel(txt Text) Label {
	return Label{txt, NewBox(txt)}
}

type Animation struct {
	Frames
	Box
}

func NewAnimation(frms Frames) Animation {
	return Animation{frms, NewBox(frms)}
}

// Filler can realize background shadow, and maybe border too.
// By introducing an image, API becomes much simpler than Web's.
// However, it is hard to adjust the size of fillers automatically
// when its parent's size changes. Nevertheless, it won't be a problem
// UI components would not change their size drastically.
type Filler struct {
	Color
	Box
}

func NewFiller(clr color.Color, extra float64) Filler {
	c := NewColor(clr)
	b := Box{
		source: c,
		Rectangle: Rectangle{
			W:      Length{extra, Extra},
			H:      Length{extra, Extra},
			Aligns: CenterMiddle,
		},
	}
	return Filler{c, b}
}
