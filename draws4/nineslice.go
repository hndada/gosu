package draws

type NineSlice struct {
	subs            [9]Sprite
	leftTopSize     XY
	rightBottomSize XY
	Box
}

func NewNineSlice(img Image, leftTopSize, rightTopSize XY) NineSlice {
	w, h := img.Size().Values()
	x0, x1, x2, x3 := 0, int(leftTopSize.X), int(w-rightTopSize.X), int(w)
	y0, y1, y2, y3 := 0, int(leftTopSize.Y), int(h-rightTopSize.Y), int(h)
	return NineSlice{
		subs: [9]Sprite{
			NewSprite(img.SubImage(x0, y0, x1, y1)),
			NewSprite(img.SubImage(x1, y0, x2, y1)),
			NewSprite(img.SubImage(x2, y0, x3, y1)),

			NewSprite(img.SubImage(x0, y1, x1, y2)),
			NewSprite(img.SubImage(x1, y1, x2, y2)),
			NewSprite(img.SubImage(x2, y1, x3, y2)),

			NewSprite(img.SubImage(x0, y2, x1, y3)),
			NewSprite(img.SubImage(x1, y2, x2, y3)),
			NewSprite(img.SubImage(x2, y2, x3, y3)),
		},
		leftTopSize:     leftTopSize,
		rightBottomSize: rightTopSize,
	}
}

func NewSimpleNineSlice(img Image, thickness float64) NineSlice {
	w, h := img.Size().Values()
	leftTopSize := XY{thickness, thickness}
	rightBottomSize := XY{w - thickness, h - thickness}
	return NewNineSlice(img, leftTopSize, rightBottomSize)
}

func (ns *NineSlice) SetSize(w, h float64) {
	// center size
	cw := w - ns.leftTopSize.X - ns.rightBottomSize.X
	ch := h - ns.leftTopSize.Y - ns.rightBottomSize.Y

	// second column
	for i := 1; i < 9; i += 3 {
		s := ns.subs[i]
		ns.subs[i]
		ns.subs[i].SetSize(cw, s.Box.Size().Y)
		ns.imgs[i].Scale(w-ns.leftTopSize.X-ns.rightBottomSize.X, ns.imgs[i].Size().Y)
	}
	// second row
	for i := 3; i < 6; i++ {
		ns.imgs[i].Scale(ns.imgs[i].Size().X, h-ns.leftTopSize.Y-ns.rightBottomSize.Y)
	}
}
