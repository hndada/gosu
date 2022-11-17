package scene

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	defaultskin "github.com/hndada/gosu/skin"
)

type LoadSkinMode int

const (
	LoadSkinDefault LoadSkinMode = iota
	LoadSkinUser
	LoadSkinPlay // refreshes every play
)

// Number sprites are used in Play and Result.
const (
	NumberDot = iota
	NumberComma
	NumberPercent
)

const (
	CursorBase = iota
	CursorAdditive
	CursorTrail
)

// In narrow meaning, Skin stands for a set of Sprites.
// In wide meaning, Skin also includes a set of sounds.
// Package defaultskin has a set of sounds.
type skinType struct {
	DefaultBackground draws.Sprite
	Number1           [13]draws.Sprite // number and sign(. , %)
	Number2           [10]draws.Sprite // number only
	Cursor            [3]draws.Sprite
}

// Unexported struct with exported function yields read-only feature.
var (
	defaultSkin skinType
	Skin        skinType
)

func init() { LoadSkin(defaultskin.FS, LoadSkinDefault) }
func LoadSkin(fsys fs.FS, mode LoadSkinMode) {
	skin := &Skin
	if mode == LoadSkinDefault {
		skin = &defaultSkin
	}
	skin.DefaultBackground = NewBackground(fsys, "interface/default-bg.jpg")
	for kind := 1; kind <= 2; kind++ {
		for i := 0; i < 10; i++ {
			s := draws.NewSprite(fsys, fmt.Sprintf("interface/number%d/%d.png", kind, i))
			s.ApplyScale(Settings.NumberScale)
			// Need to set same base line, since each number has different height.
			if i == 0 {
				s.Locate(ScreenSizeX, 0, draws.RightTop)
			} else {
				s.Locate(ScreenSizeX, skin.Number1[0].H()-s.H(), draws.RightTop)
			}
			switch kind {
			case 1:
				skin.Number1[i] = s
			case 2:
				skin.Number2[i] = s
			}
		}
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("number/score/%s.png", name))
		s.ApplyScale(Settings.NumberScale)
		// Cursor should be at CenterMiddle in circle mode (in far future)
		s.Locate(ScreenSizeX, 0, draws.RightTop)
		skin.Number1[10+i] = s
	}
	for i, name := range []string{"base", "additive", "trail"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.ApplyScale(Settings.CursorScale)
		s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		skin.Cursor[i] = s
	}
	skin.fillBlank([]skinType{{}, defaultSkin, Skin}[mode])
}

func (skin *skinType) fillBlank(base skinType) {
	if !skin.DefaultBackground.IsValid() {
		skin.DefaultBackground = base.DefaultBackground
	}
	for _, s := range skin.Number1 {
		if !s.IsValid() {
			skin.Number1 = base.Number1
			break
		}
	}
	for _, s := range skin.Number2 {
		if !s.IsValid() {
			skin.Number2 = base.Number2
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
func NewBackground(fsys fs.FS, name string) draws.Sprite {
	s := draws.NewSprite(fsys, name)
	s.ApplyScale(ScreenSizeX / s.W())
	s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	return s
}
