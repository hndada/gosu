package draws

type NineSlice struct {
	subs   [9]Sprite
	w0, h0 float64 // left top
	w2, h2 float64 // right bottom
	Box
}

func NewNineSlice(img Image, w0, h0, w2, h2 float64) NineSlice {
	w, h := img.Size().Values()
	x0, x1, x2, x3 := 0, int(w0), int(w-w2), int(w)
	y0, y1, y2, y3 := 0, int(h0), int(h-h2), int(h)
	ns := NineSlice{
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
		w0: w0, h0: h0,
		w2: w2, h2: h2,
		Box: NewBox(img),
	}

	for i := 0; i < 9; i++ {
		ns.subs[i].SetOrigin(&ns.Box)
	}
	ns.SetSize(w, h)
	return ns
}

func NewSimpleNineSlice(img Image, thickness float64) NineSlice {
	return NewNineSlice(img, thickness, thickness, thickness, thickness)
}

func (ns *NineSlice) SetSize(w, h float64) {
	w0, h0 := ns.w0, ns.h0
	w1 := w - ns.w0 - ns.w2
	h1 := h - ns.h0 - ns.h2
	w2, h2 := ns.w2, ns.h2

	for i := 0; i < 9; i++ {
		var w, h float64
		var x, y float64
		switch i / 3 { // row
		case 0:
			h, y = h0, 0
		case 1:
			h, y = h1, h0
		case 2:
			h, y = h2, h1+h0
		}
		switch i % 3 { // column
		case 0:
			w, x = w0, 0
		case 1:
			w, x = w1, w0
		case 2:
			w, x = w2, w1+w0
		}
		ns.subs[i].Box.Size.SetValues(w, h)
		ns.subs[i].Box.Position.SetValues(x, y)
	}
}

func (ns NineSlice) Draw(dst Image) {
	for i := 0; i < 9; i++ {
		ns.subs[i].Draw(dst)
	}
}
