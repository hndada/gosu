package draws

import "image/color"

// Filler can realize background shadow, and maybe border too.
// By introducing an image, API becomes much simpler than Web's.
// However, it is hard to adjust the size of fillers automatically
// when its parent's size changes. Nevertheless, it won't be a problem
// UI components would not change their size drastically.

func NewFiller(base *Box, clr color.Color, extra float64) Filler {
	return Filler{
		Color: NewColor(clr),
		Box: Box{
			Base:   base,
			Size:   NewLength2(&base.Size, extra, extra, Extra),
			Aligns: CenterMiddle,
		},
	}
}
