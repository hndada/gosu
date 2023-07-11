package draws

import (
	"io/fs"
)

type Animation []Sprite

func NewAnimation(srcs any) Animation {
	switch srcs := srcs.(type) {
	case []Image:
		return newAnimationFromImages(srcs)
	}
	return nil
}
func newAnimationFromImages(imgs []Image) Animation {
	a := make(Animation, len(imgs))
	for i, img := range imgs {
		a[i] = NewSprite(img)
	}
	return a
}

func NewAnimationFromFile(fsys fs.FS, name string) Animation {
	return NewAnimation(NewImagesFromFile(fsys, name))
}

func (a Animation) IsEmpty() bool {
	return len(a) <= 1 && a[0].IsEmpty()
}
