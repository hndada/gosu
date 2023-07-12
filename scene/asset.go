package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
)

const (
	ScreenSizeX = 1600
	ScreenSizeY = 900
)

const (
	CursorBase = iota
	CursorAdditive
	CursorTrail
)

// Asset is previously known as "Skin"
// Skin might not be clear.
// Assets would be confusing name for singleton.
type Asset struct {
	Cursor            [3]draws.Sprite
	DefaultBackground draws.Sprite
	BoxMask           draws.Sprite
	Clear             draws.Sprite
	// Intro   draws.Sprite
	// Loading draws.Sprite

	Enter      audios.Sound
	Swipe      audios.SoundBag
	Tap        audios.SoundBag
	Toggle     [2]audios.Sound
	Transition [2]audios.Sound
}

var TheAsset = Asset{}

// Wrapping with each block looks good in terms of readability and reliability.
func LoadTheAsset(fsys fs.FS) {
	// Cursor should be at CenterMiddle in circle mode (in far future)
	for i, name := range []string{"base", "additive", "trail"} {
		s := draws.NewSpriteFromFile(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.MultiplyScale(TheSettings.CursorScale)
		s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.LeftTop)
		TheAsset.Cursor[i] = s
	}
	TheAsset.DefaultBackground = NewBackgroundFromFile(fsys, "interface/default-bg.jpg")
	TheAsset.BoxMask = draws.NewSpriteFromFile(fsys, "interface/box-mask.png")
	{
		s := draws.NewSpriteFromFile(fsys, "interface/clear.png")
		s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		s.MultiplyScale(TheSettings.ClearScale)
		TheAsset.Clear = s
	}
	{
		streamer, _, _ := audios.DecodeFromFile(fsys, "sound/ringtone2_loop.wav")
		TheAsset.Enter = streamer
	}
	TheAsset.Swipe = audios.NewSoundBag(fsys, "sound/swipe.wav")
	TheAsset.Tap = audios.NewSoundBag(fsys, "sound/tap.wav")
	for i, name := range []string{"off", "on"} {
		name := fmt.Sprintf("sound/toggle/%s.wav", name)
		streamer, _, _ := audios.DecodeFromFile(fsys, name)
		TheAsset.Toggle[i] = streamer
	}
	for i, name := range []string{"down", "up"} {
		name := fmt.Sprintf("sound/transition/%s.wav", name)
		streamer, _, _ := audios.DecodeFromFile(fsys, name)
		TheAsset.Transition[i] = streamer
	}
}
