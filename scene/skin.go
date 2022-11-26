package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
)

type Skin struct {
	Type              int
	DefaultBackground draws.Sprite
	Cursor            [3]draws.Sprite
	BoxMask           draws.Sprite

	Enter      audios.Sound
	Swipe      audios.SoundBag
	Tap        audios.SoundBag
	Toggle     [2]audios.Sound
	Transition [2]audios.Sound
}

const (
	CursorBase = iota
	CursorAdditive
	CursorTrail
)

var (
	DefaultSkin = Skin{Type: mode.Default}
	UserSkin    = Skin{Type: mode.User}
)

func (skin *Skin) Load(fsys fs.FS) {
	skin.DefaultBackground = mode.UserSkin.DefaultBackground
	for i, name := range []string{"base", "additive", "trail"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.ApplyScale(S.CursorScale)
		// Cursor should be at CenterMiddle in circle mode (in far future)
		s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		skin.Cursor[i] = s
	}
	skin.BoxMask = draws.NewSprite(fsys, "interface/box-mask.png")
	skin.Enter = audios.NewSound(fsys, "sound/ringtone2_loop.wav")
	skin.Swipe = audios.NewSoundBag(fsys, "sound/swipe")
	skin.Tap = audios.NewSoundBag(fsys, "sound/tap")
	for i, name := range []string{"off", "on"} {
		name := fmt.Sprintf("sound/toggle/%s.wav", name)
		skin.Toggle[i] = audios.NewSound(fsys, name)
	}
	for i, name := range []string{"off", "on"} {
		name := fmt.Sprintf("sound/transition/%s.wav", name)
		skin.Transition[i] = audios.NewSound(fsys, name)
	}
	base := []Skin{{}, DefaultSkin, UserSkin}[skin.Type]
	skin.fillBlank(base)
}
func (skin *Skin) fillBlank(base Skin) {
	for _, s := range skin.Cursor {
		if !s.IsValid() {
			skin.Cursor = base.Cursor
			break
		}
	}
	if !skin.Enter.IsValid() {
		skin.Enter = base.Enter
	}
	for _, s := range skin.Swipe {
		if !s.IsValid() {
			skin.Swipe = base.Swipe
			break
		}
	}
	for _, s := range skin.Tap {
		if !s.IsValid() {
			skin.Tap = base.Tap
			break
		}
	}
	for _, s := range skin.Toggle {
		if !s.IsValid() {
			skin.Toggle = base.Toggle
			break
		}
	}
	for _, s := range skin.Transition {
		if !s.IsValid() {
			skin.Transition = base.Transition
			break
		}
	}
}
