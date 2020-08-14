package ebitenui

import "github.com/hajimehoshi/ebiten"

// AttachH attaches input images horizontally.
// AttachH supposes all input images have proper height.
func AttachH(imgs ...*ebiten.Image) *ebiten.Image {
	var w, h int
	for _, i := range imgs {
		x, y := i.Size()
		if h < y {
			h = y
		}
		w += x
	}
	img, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	op := &ebiten.DrawImageOptions{}
	for _, i := range imgs {
		img.DrawImage(i, op)
		x, _ := i.Size()
		op.GeoM.Translate(float64(x), 0)
	}
	return img
}

// AttachV attaches input images vertically.
// AttachV supposes all input images have proper width.
func AttachV(imgs ...*ebiten.Image) *ebiten.Image {
	var w, h int
	for _, i := range imgs {
		x, y := i.Size()
		if x > w {
			w = x
		}
		h += y
	}
	img, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	op := &ebiten.DrawImageOptions{}
	for _, i := range imgs {
		img.DrawImage(i, op)
		_, y := i.Size()
		op.GeoM.Translate(0, float64(y))
	}
	return img
}
