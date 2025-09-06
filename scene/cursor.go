package scene

import (
	"github.com/hndada/gosu/draws"
)

type Cursor struct {
	base     draws.Sprite
	additive draws.Sprite
	trail    draws.Sprite
}

// Cursor should be at CenterMiddle in circle mode (in far future)
func NewCursor(res *Resources, opts *Options) Cursor {
	cursor := Cursor{}
	{
		s := draws.NewSprite(res.CursorBaseImage)
		s.Scale(opts.MouseCursorImageScale)
		cursor.base = s
	}
	{
		s := draws.NewSprite(res.CursorAdditiveImage)
		s.Scale(opts.MouseCursorImageScale)
		cursor.additive = s
	}
	{
		s := draws.NewSprite(res.CursorTrailImage)
		s.Scale(opts.MouseCursorImageScale)
		cursor.trail = s
	}
	return cursor
}
