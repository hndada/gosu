package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/audios"
)

type Sounds struct {
	Select     audios.Sound
	Swipe      audios.SoundBag
	Tap        audios.SoundBag
	Toggle     [2]audios.Sound
	Transition [2]audios.Sound
}

// Todo: need a test whether fsys is immutable
func loadSounds(fsys fs.FS, mode LoadSkinMode) {
	const prefix = "interface/sound/"
	var sounds Sounds
	sounds.Select = audios.NewSound(fsys, prefix+"ringtone2_loop.wav")
	sounds.Swipe = audios.NewSoundBag(fsys, prefix+"swipe")
	sounds.Tap = audios.NewSoundBag(fsys, prefix+"tap")
	for i, name := range []string{"off", "on"} {
		name := fmt.Sprintf("%s/toggle/%s.wav", prefix, name)
		sounds.Toggle[i] = audios.NewSound(fsys, name)
	}
	for i, name := range []string{"off", "on"} {
		name := fmt.Sprintf("%s/transition/%s.wav", prefix, name)
		sounds.Transition[i] = audios.NewSound(fsys, name)
	}
	base := []Sounds{{}, defaultSkin.Sounds, userSkin.Sounds}[mode]
	sounds.loadBase(base)
}
func (sounds *Sounds) loadBase(base Sounds) {
	if !sounds.Select.IsValid() {
		sounds.Select = base.Select
	}
	for _, s := range sounds.Swipe {
		if !s.IsValid() {
			sounds.Swipe = base.Swipe
			break
		}
	}
	for _, s := range sounds.Tap {
		if !s.IsValid() {
			sounds.Tap = base.Tap
			break
		}
	}
	for _, s := range sounds.Toggle {
		if !s.IsValid() {
			sounds.Toggle = base.Toggle
			break
		}
	}
	for _, s := range sounds.Transition {
		if !s.IsValid() {
			sounds.Transition = base.Transition
			break
		}
	}
}
