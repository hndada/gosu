package draws

import "github.com/hajimehoshi/ebiten/v2"

type Animation []Sprite

func NewAnimation(path string) Animation {
	return NewAnimationFromImages(NewImages(path))
}
func NewAnimationFromImages(images []*ebiten.Image) (a Animation) {
	a = make(Animation, len(images))
	for i, image := range images {
		a[i] = NewSpriteFromImage(image)
	}
	return
}
