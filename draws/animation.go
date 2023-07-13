package draws

import (
	"io/fs"
)

type Animation []Sprite

func NewAnimation(srcs any) Animation {
	switch srcs := srcs.(type) {
	case Frames:
		return newAnimationFromFrames(srcs)
	}
	return nil
}

func newAnimationFromFrames(frames Frames) Animation {
	a := make(Animation, len(frames))
	for i, img := range frames {
		a[i] = NewSprite(img)
	}
	return a
}

func NewAnimationFromFile(fsys fs.FS, name string) Animation {
	return NewAnimation(NewFramesFromFilename(fsys, name))
}

func (a Animation) IsEmpty() bool {
	return len(a) <= 1 && a[0].IsEmpty()
}
