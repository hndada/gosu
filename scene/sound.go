package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/audios"
)

type Sounds struct {
	empty      bool
	Select     audios.Sound
	Swipe      audios.Sound
	Tap        audios.Sound
	Toggle     [2]audios.Sound
	Transition [2]audios.Sound
}

var (
	defaultSounds Sounds
	currentSounds Sounds
)

// Todo: need a test whether fsys is immutable
func loadSound(fsys fs.FS, base Sounds) {
	sounds := &currentSounds
	if base.empty {
		sounds = &defaultSounds
	}
	for i, name := range []string{"tap/0", "old/restart", "swipe"} {
		sound, err := audios.NewSound(fsys, fmt.Sprintf("sound/%s.wav", name))
		if err != nil {
			if base.empty {
				panic("fail to load default sounds")
			}
			sound = []audios.Sound{base.Tap, base.Select, base.Swipe}[i]
		}
		switch i {
		case 0:
			sounds.Tap = sound
		case 1:
			sounds.Select = sound
		case 2:
			sounds.Swipe = sound
		}
	}
	for i, name := range []string{"off", "on"} {
		name := fmt.Sprintf("sound/toggle/%s.wav", name)
		sound, err := audios.NewSound(fsys, name)
		if err != nil {
			if base.empty {
				panic("fail to load default sounds")
			}
			sound = base.Toggle[i]
		}
		sounds.Toggle[i] = sound
	}
	for i, name := range []string{"down", "up"} {
		name := fmt.Sprintf("sound/transition/%s.wav", name)
		sound, err := audios.NewSound(fsys, name)
		if err != nil {
			if base.empty {
				panic("fail to load default sounds")
			}
			sound = base.Transition[i]
		}
		sounds.Transition[i] = sound
	}
}
