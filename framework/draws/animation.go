package draws

type Animation []Sprite

func NewAnimation(path string) Animation {
	return NewAnimationFromImages(LoadImages(path))
}
func NewAnimationFromImages(images []Image) (a Animation) {
	a = make(Animation, len(images))
	for i, image := range images {
		a[i] = NewSpriteFromSource(image)
	}
	return
}
