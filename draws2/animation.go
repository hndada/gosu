package draws

type Animation struct {
	Frames
	Options
}

func NewAnimation(src Frames) Animation {
	return Animation{src, NewOptions(src)}
}

func (a Animation) Draw(dst Image) {
	if a.IsEmpty() {
		return
	}
	src := a.Frame().Image
	dst.DrawImage(src, a.imageOp())
}
