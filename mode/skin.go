package mode

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/render"
)

// ScreenSize is a logical size of in-game screen.
const (
	screenSizeX = 1600
	screenSizeY = 900

	ScreenSizeX = screenSizeX
	ScreenSizeY = screenSizeY
)

// Skin is a set of Sprites.
// mode.Skin is a general skin for all modes.
var (
	DefaultBackground render.Sprite
	CursorSprites     [2]render.Sprite // 0: cursor // 1: additive cursor
	// CursorTailSprite   Sprite
)

func LoadSkin() {
	DefaultBackground = render.Sprite{
		I:      render.NewImage("skin/default-bg.jpg"),
		Filter: ebiten.FilterLinear,
	}
	DefaultBackground.SetWidth(screenSizeX)
	DefaultBackground.SetCenterY(screenSizeY / 2)

	for i, name := range []string{"menu-cursor.png", "menu-cursor-additive.png"} {
		s := render.Sprite{
			I:      render.NewImage(fmt.Sprintf("skin/cursor/%s", name)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(CursorScale)
		CursorSprites[i] = s
	}
}
