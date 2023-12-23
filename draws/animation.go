package draws

import (
	"io/fs"
)

type Animation struct {
	Sprites []Sprite
	tick    int
	maxTick int
}

func NewAnimation(srcs any, maxTick int) Animation {
	switch srcs := srcs.(type) {
	case []Sprite:
		return Animation{Sprites: srcs, maxTick: maxTick}
	case ImageSequence:
		return newAnimationFromImageSequence(srcs, maxTick)
	}
	return Animation{}
}

func newAnimationFromImageSequence(seq ImageSequence, maxTick int) (a Animation) {
	a.Sprites = make([]Sprite, len(seq))
	for i, img := range seq {
		a.Sprites[i] = NewSprite(img)
	}
	return a
}

func NewAnimationFromFile(fsys fs.FS, name string, maxTick int) Animation {
	return NewAnimation(NewImageSequenceFromFilename(fsys, name), maxTick)
}

func (a *Animation) MultiplyScale(scale float64) {
	for i := range a.Sprites {
		a.Sprites[i].MultiplyScale(scale)
	}
}

func (a *Animation) Locate(x, y float64, anchor Anchor) {
	for i := range a.Sprites {
		a.Sprites[i].Locate(x, y, anchor)
	}
}

func (a *Animation) Tick() {
	if a.maxTick > 0 {
		a.tick = (a.tick % a.maxTick) + 1
	}
}

func (a Animation) Frame() Sprite {
	if len(a.Sprites) == 0 {
		return Sprite{}
	}
	if a.maxTick == 0 {
		return a.Sprites[0]
	}
	progress := float64(a.tick%a.maxTick) / float64(a.maxTick)
	count := float64(len(a.Sprites))
	index := int(progress * count)
	return a.Sprites[index]
}

func (a Animation) Draw(screen Image, op Op) {
	a.Frame().Draw(screen, op)
}

func (a Animation) IsEmpty() bool {
	return len(a.Sprites) <= 1 && a.Sprites[0].Source.IsEmpty()
}

func (a *Animation) Reset() { a.tick = 0 }
