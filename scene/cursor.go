package scene

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
		s := draws.NewSprite(res.CursorBase)
		s.Scale(opts.MouseCursorImageScale)
		cmp.base = s
	}
	{
		s := draws.NewSprite(res.CursorAdditive)
		s.Scale(opts.MouseCursorImageScale)
		cmp.additive = s
	}
	{
		s := draws.NewSprite(res.CursorTrail)
		s.Scale(opts.MouseCursorImageScale)
		cmp.trail = s
	}
	return
}
