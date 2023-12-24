package draws

// Since Box is initialized with given Image,
// Source such as image or text should be unexported.
type Sprite struct {
	img Image
	Box
}

func NewSprite(img Image) Sprite {
	return Sprite{
		img: img,
		Box: NewBox(img),
	}
}

func (s Sprite) Draw(dst Image) {
	s.Box.Draw(dst, s.img)
}
