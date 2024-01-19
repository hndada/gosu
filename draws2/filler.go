package draws

import "image/color"

// Filler can realize background shadow, and maybe border too.
// By introducing an image, API becomes much simpler than Web's.
// However, it is hard to adjust the size of fillers automatically
// when its parent's size changes. Nevertheless, it won't be a problem
// UI components would not change their size drastically.
type Filler struct {
	Color
	Options
}

func NewFiller(base *Rectangle, clr color.Color, extra float64) Filler {
	return Filler{
		Color: NewColor(clr),
		Options: Options{
			Rectangle: Rectangle{
				Base:   base,
				Size:   NewLength2(&base.Size, extra, extra, Extra),
				Aligns: CenterMiddle,
			},
		},
	}
}

func (f Filler) Draw(dst Image) {
	if f.IsEmpty() {
		return
	}
	sub := dst.Sub(f.Min(), f.Max())
	src := f.Color.Color
	sub.Fill(src)
}
