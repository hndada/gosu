package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type CursorResources struct {
	base     draws.Image
	additive draws.Image
	trail    draws.Image
}

func (res *CursorResources) Load(fsys fs.FS) {
	const (
		base = iota
		additive
		trail
	)
	for i, name := range []string{"base", "additive", "trail"} {
		fname := fmt.Sprintf("interface/cursor/%s.png", name)
		switch i {
		case base:
			res.base = draws.NewImageFromFile(fsys, fname)
		case additive:
			res.additive = draws.NewImageFromFile(fsys, fname)
		case trail:
			res.trail = draws.NewImageFromFile(fsys, fname)
		}
	}
}

type CursorOptions struct {
	Scale float64
}

type CursorComponent struct {
	base     draws.Sprite
	additive draws.Sprite
	trail    draws.Sprite
}

// Cursor should be at CenterMiddle in circle mode (in far future)
func NewCursorComponent(res CursorResources, opts CursorOptions) (cmp CursorComponent) {
	{
		s := draws.NewSprite(res.base)
		s.MultiplyScale(opts.Scale)
		cmp.base = s
	}
	{
		s := draws.NewSprite(res.additive)
		s.MultiplyScale(opts.Scale)
		cmp.additive = s
	}
	{
		s := draws.NewSprite(res.trail)
		s.MultiplyScale(opts.Scale)
		cmp.trail = s
	}
	return
}
