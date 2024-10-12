package game

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/plays/piano"
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
	BoxMaskImage           draws.Image
	CursorBaseImage        draws.Image
	CursorAdditiveImage    draws.Image
	CursorTrailImage       draws.Image

	Piano *piano.Resources
}

func NewResources(fsys fs.FS) (res *Resources) {
	res = &Resources{}
	{
		fname := "interface/default-bg.jpg"
		res.DefaultBackgroundImage = draws.NewImageFromFile(fsys, fname)
	}
	{
		fname := "interface/box-mask.png"
		res.BoxMaskImage = draws.NewImageFromFile(fsys, fname)
	}
	{
		fname := "interface/cursor/base.png"
		res.CursorBaseImage = draws.NewImageFromFile(fsys, fname)
	}
	{
		fname := "interface/cursor/additive.png"
		res.CursorAdditiveImage = draws.NewImageFromFile(fsys, fname)
	}
	{
		fname := "interface/cursor/trail.png"
		res.CursorTrailImage = draws.NewImageFromFile(fsys, fname)
	}
	res.Piano = piano.NewResources(fsys)
	return
}
