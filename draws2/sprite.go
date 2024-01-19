package draws

// Image does not have Draw method.
// Sprite and Animation have.
type Sprite struct {
	Image
	Options
}

func NewSprite(src Image) Sprite {
	return Sprite{src, NewOptions(src)}
}

func (s Sprite) Draw(dst Image) {
	if s.IsEmpty() {
		return
	}
	src := s.Image.Image
	dst.DrawImage(src, s.imageOp())
}
