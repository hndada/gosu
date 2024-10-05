package scene

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game/piano"
)

// Todo: deal with two kinds of values: file path and directory path.
const (
	SoundToggleOff      = "interface/sound/toggle/off.wav"
	SoundToggleOn       = "interface/sound/toggle/on.wav"
	SoundTransitionDown = "interface/sound/transition/down.wav"
	SoundTransitionUp   = "interface/sound/transition/up.wav"
	SoundTaps           = "interface/sound/tap/"
	SoundSwipes         = "interface/sound/swipe/"
)

// SearchBoxSprite's Color: RGBA{128, 128, 128, 128}
// SearchBoxSprite's Height: 25
type Resources struct {
	DefaultBackgroundImage draws.Image
	CursorBase             draws.Image
	CursorAdditive         draws.Image
	CursorTrail            draws.Image

	Piano *piano.Resources
}

func NewResources(fsys fs.FS) (res *Resources) {
	res = &Resources{}
	{
		fname := "interface/default-bg.jpg"
		res.DefaultBackgroundImage = draws.NewImageFromFile(fsys, fname)
	}

	{
		fname := "interface/cursor/base.png"
		res.CursorBase = draws.NewImageFromFile(fsys, fname)
	}
	{
		fname := "interface/cursor/additive.png"
		res.CursorAdditive = draws.NewImageFromFile(fsys, fname)
	}
	{
		fname := "interface/cursor/trail.png"
		res.CursorTrail = draws.NewImageFromFile(fsys, fname)
	}
	res.Piano = piano.NewResources(fsys)
	return
}
