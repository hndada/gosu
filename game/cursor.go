package game

import (
	"github.com/hndada/gosu/draws"
)

type CursorComponent struct {
	base     draws.Sprite
	additive draws.Sprite
	trail    draws.Sprite
}

// Cursor should be at CenterMiddle in circle mode (in far future)
func NewCursorComponent(res *Resources, opts *Options) (cmp CursorComponent) {
	{
		s := draws.NewSprite(res.CursorBaseImage)
		s.Scale(opts.MouseCursorImageScale)
		cmp.base = s
	}
	{
		s := draws.NewSprite(res.CursorAdditiveImage)
		s.Scale(opts.MouseCursorImageScale)
		cmp.additive = s
	}
	{
		s := draws.NewSprite(res.CursorTrailImage)
		s.Scale(opts.MouseCursorImageScale)
		cmp.trail = s
	}
	return
}
