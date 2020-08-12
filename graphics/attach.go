package graphics

import "github.com/hajimehoshi/ebiten"

// Suppose all images have proper height e.g., same
func AttachH(imgs ...*ebiten.Image) *ebiten.Image { // horizontal
	var w, h int
	for _, i := range imgs {
		x, y := i.Size()
		if y > h {
			y = h
		}
		w += x
	}
	img, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	op:= &ebiten.DrawImageOptions{}
	for _, i := range imgs {
		img.DrawImage(i, op)
		x, _ := i.Size()
		op.GeoM.Translate(float64(x), 0)
	}
	return img
}

func AttachV(imgs ...*ebiten.Image) *ebiten.Image { // vertical
	var w, h int
	for _, i := range imgs {
		x, y := i.Size()
		if x > w {
			w = x
		}
		h += y
	}
	img, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	op:= &ebiten.DrawImageOptions{}
	for _, i := range imgs {
		img.DrawImage(i, op)
		_, y := i.Size()
		op.GeoM.Translate(0, float64(y))
	}
	return img
}
