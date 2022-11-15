package scene

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

// Skin is a set of Sprites and sounds.
type Skin struct {
	empty             bool
	DefaultBackground draws.Sprite
	Cursor            [3]draws.Sprite
	Score             [10]draws.Sprite
	Sign              [3]draws.Sprite
	// ChartItemBoxSprite draws.Sprite
	// ChartLevelBoxSprite draws.Sprite
	Sounds
}

// Unexported struct with exported function yields read-only feature.
var (
	defaultSkin Skin
	currentSkin Skin
)

func initSkin() {
	//go:embed skin/*
	var fs embed.FS
	LoadSkin(fs, Skin{empty: true})
}
func DefaultSkin() Skin { return defaultSkin }
func CurrentSkin() Skin { return currentSkin }

// Todo: skip when not existed
func LoadSkin(fsys fs.FS, base Skin) {
	skin := &currentSkin
	if base.empty {
		skin = &defaultSkin
	}
	settings := currentSettings
	{
		sprite := NewBackground(fsys, "default-bg.jpg")
		if !sprite.IsValid() {
			if base.empty {
				panic("fail to load default skin")
			}
			sprite = base.DefaultBackground
		}
		skin.DefaultBackground = sprite
	}
	for i, name := range []string{"base", "additive", "trail"} {
		s := draws.NewSprite(fsys, fmt.Sprintf("cursor/%s.png", name))
		s.ApplyScale(settings.CursorScale)
		s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
		skin.Cursor[i] = s
	}
	loadSound(fsys, base.Sounds)
}
func NewBackground(fsys fs.FS, name string) draws.Sprite {
	s := draws.NewSprite(fsys, name)
	s.ApplyScale(ScreenSizeX / s.W())
	s.Locate(ScreenSizeX/2, ScreenSizeY/2, draws.CenterMiddle)
	return s
}

// {
// 	s := draws.NewSprite(fsys, "box-mask.png")
// 	s.SetSize(ChartInfoBoxWidth, ChartInfoBoxHeight)
// 	s.Locate(ScreenSizeX+chartInfoBoxshrink, ScreenSizeY/2, draws.RightMiddle)
// 	ChartItemBoxSprite = s
// }
