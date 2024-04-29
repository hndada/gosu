package draws

type Sprite struct {
	Source Image
	Box
}

func NewSprite(img Image) Sprite {
	return Sprite{
		Source: img,
		Box:    NewBox(img),
	}
}

// sub.Fill might fill the destination image permanently.
func (s Sprite) Draw(dst Image) {
	if s.Source.IsEmpty() {
		return
	}
	dst.DrawImage(s.Source.Image, s.op())
}
