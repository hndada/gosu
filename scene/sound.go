package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/audios"
	defaultskin "github.com/hndada/gosu/skin"
)

type Sounds struct {
	Enter      audios.Sound
	Swipe      audios.SoundBag
	Tap        audios.SoundBag
	Toggle     [2]audios.Sound
	Transition [2]audios.Sound
}

var (
	defaultSounds Sounds
	userSounds    Sounds
)

func init()                 { LoadSounds(defaultskin.FS, LoadSkinDefault) }
func DefaultSounds() Sounds { return defaultSounds }
func UserSounds() Sounds    { return userSounds }
func LoadSounds(fsys fs.FS, mode LoadSkinMode) {
	const prefix = "sound/"
	var sounds Sounds
	sounds.Enter = audios.NewSound(fsys, prefix+"ringtone2_loop.wav")
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
	base := []Sounds{{}, defaultSounds, userSounds}[mode]
	sounds.fillBlank(base)
}
func (sounds *Sounds) fillBlank(base Sounds) {
	if !sounds.Enter.IsValid() {
		sounds.Enter = base.Enter
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
