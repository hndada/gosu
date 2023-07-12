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
	Swipe      audios.SoundPod
	Tap        audios.SoundPod
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
	var SoundVolume float64
	{
		sound, err := audios.NewSound(fsys, "sound/ringtone2_loop.wav", &SoundVolume)
		if err != nil {
			panic(err)
		}
		TheAsset.Enter = sound
	}
	{
		subFS, err := fs.Sub(fsys, "sound/swipe")
		if err != nil {
			panic(err)
		}
		TheAsset.Swipe = audios.NewSoundPod(subFS, &SoundVolume)
	}
	{
		subFS, err := fs.Sub(fsys, "sound/tap")
		if err != nil {
			panic(err)
		}
		TheAsset.Tap = audios.NewSoundPod(subFS, &SoundVolume)
	}
	for i, name := range []string{"off", "on"} {
		name := fmt.Sprintf("sound/toggle/%s.wav", name)
		sound, err := audios.NewSound(fsys, name, &SoundVolume)
		if err != nil {
			panic(err)
		}
		TheAsset.Toggle[i] = sound
	}
	for i, name := range []string{"down", "up"} {
		name := fmt.Sprintf("sound/transition/%s.wav", name)
		sound, err := audios.NewSound(fsys, name, &SoundVolume)
		if err != nil {
			panic(err)
		}
		TheAsset.Transition[i] = sound
	}
}
