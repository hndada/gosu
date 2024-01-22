package draws

type Sprite struct {
	Image
	Box
}

func NewSprite(img Image) Sprite {
	return Sprite{
		Image: img,
		Box:   NewBox(img),
	}
}

// sub.Fill might fill the destination image permanently.
func (s Sprite) Draw(dst Image) {
	if s.IsEmpty() {
		return
	}
	dst.DrawImage(s.Image.Image, s.op())
}
