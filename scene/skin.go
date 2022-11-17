package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

type LoadSkinMode int

const (
	LoadSkinDefault = iota
	LoadSkinUser
	LoadSkinPlay
)

// Skin is a set of Sprites and sounds.
type Skin struct {
	DefaultBackground draws.Sprite
	Score             [13]draws.Sprite // number + sign(. , %)
	Combo             [10]draws.Sprite // number only
	Cursor            [3]draws.Sprite
	Sounds
}

// Unexported struct with exported function yields read-only feature.
var (
	defaultSkin Skin
	userSkin    Skin
	// playSkin    Skin // refreshes every play
)

func initSkin() {
	LoadSkin(defaultskin.FS, LoadSkinDefault)
}
func DefaultSkin() Skin { return defaultSkin }
func CurrentSkin() Skin { return userSkin }

// LoadSkin either assigns data to skin or returns skin
// based on the mode.
// Todo: skip when not existed
func LoadSkin(fsys fs.FS, mode LoadSkinMode) {
	var skin Skin
	skin.DefaultBackground = NewBackground(fsys, "interface/default-bg.jpg")
	for i := 0; i < 10; i++ {
		s := draws.NewSprite(fsys, fmt.Sprintf("number/score/%d.png", i))
		s.ApplyScale(settings.ScoreScale)
		// Need to set same base line, since each number has different height.
		if i == 0 {
			s.Locate(ScreenSizeX, 0, draws.RightTop)
		} else {
			s.Locate(ScreenSizeX, skin.Score[0].H()-s.H(), draws.RightTop)
		}
		skin.Score[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("number/score/%s.png", name))
		s.ApplyScale(settings.ScoreScale)
		s.Locate(ScreenSizeX, 0, draws.RightTop)
		skin.Score[10+i] = s
	}
	for i, name := range []string{"base", "additive", "trail"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.ApplyScale(settings.CursorScale)
		s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		skin.Cursor[i] = s
	}
	loadSounds(fsys, mode)
	skin.loadBase([]Skin{{}, defaultSkin, userSkin}[mode])
	switch mode {
	case LoadSkinDefault:
		defaultSkin = skin
	case LoadSkinUser:
		userSkin = skin
		// case LoadSkinPlay:
		// 	playSkin = skin
	}
}
func NewBackground(fsys fs.FS, name string) draws.Sprite {
	s := draws.NewSprite(fsys, name)
	s.ApplyScale(ScreenSizeX / s.W())
	s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	return s
}
func (skin *Skin) loadBase(base Skin) {
	if !skin.DefaultBackground.IsValid() {
		skin.DefaultBackground = base.DefaultBackground
	}
	for _, s := range skin.Score {
		if !s.IsValid() {
			skin.Score = base.Score
			break
		}
	}
	for _, s := range skin.Combo {
		if !s.IsValid() {
			skin.Combo = base.Combo
			break
		}
	}
	for _, s := range skin.Cursor {
		if !s.IsValid() {
			skin.Cursor = base.Cursor
			break
		}
	}
}

// v := reflect.ValueOf(dst)
// t := v.Type()
// for i := 0; i < v.NumField(); i++ {
// 	if v.Field(i).Interface() == nil {}
// 	t.Field(i).Name
// }
