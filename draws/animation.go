package draws

type Animation []Sprite

func NewAnimation(path string) Animation {
	return NewAnimationFromImages(NewImages(path))
}
func NewAnimationFromImages(images []Image) (a Animation) {
	a = make(Animation, len(images))
	for i, image := range images {
		a[i] = NewSprite(image)
	}
	return
}
