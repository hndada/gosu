package gosu

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
)

// ScreenSize is a logical size of in-game screen.
const (
	screenSizeX = 1600
	screenSizeY = 900

	ScreenSizeX = screenSizeX
	ScreenSizeY = screenSizeY
)

// Skin is a set of Sprites.
var (
	DefaultBackground  draws.Sprite
	CursorSprites      [3]draws.Sprite
	ChartInfoBoxSprite draws.Sprite // Todo: various box sprite

	ScoreSprites [10]draws.Sprite
	SignSprites  [3]draws.Sprite
	// ComboSprites [10]draws.Sprite
)

const (
	CursorSpriteBase = iota
	CursorSpriteAdditive
	CursorSpriteTrail
)

func LoadGeneralSkin() {
	DefaultBackground = NewBackground("skin/default-bg.jpg")
	// Todo: cursor may have outer circle
	names := []string{"menu-cursor", "menu-cursor-additive", "cursortrail"}
	for i, name := range names {
		s := draws.NewSprite(fmt.Sprintf("skin/cursor/%s.png", name))
		s.SetScale(CursorScale)
		s.SetPosition(screenSizeX/2, screenSizeY/2, draws.OriginCenterMiddle)
		CursorSprites[i] = s
	}
	{
		s := draws.NewSprite("skin/box.png")
		scaleW := ChartInfoBoxWidth / s.W()
		scaleH := ChartInfoBoxHeight / s.H()
		s.SetScaleXY(scaleW, scaleH, ebiten.FilterLinear)
		s.SetPosition(screenSizeX+chartInfoBoxshrink, screenSizeY/2, draws.OriginRightMiddle)
		ChartInfoBoxSprite = s
	}
	for i := 0; i < 10; i++ {
		s := draws.NewSprite(fmt.Sprintf("skin/score/%d.png", i))
		s.SetScale(ScoreScale)
		if i == 0 {
			s.SetPosition(screenSizeX, 0, draws.OriginRightTop)
		} else { // Need to set same base line, since each number has different height.
			s.SetPosition(screenSizeX, ScoreSprites[0].H()-s.H(), draws.OriginRightTop)
		}
		ScoreSprites[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.NewSprite(fmt.Sprintf("skin/score/%s.png", name))
		s.SetScale(ScoreScale)
		s.SetPosition(screenSizeX, 0, draws.OriginRightTop)
		SignSprites[i] = s
	}
}
func NewBackground(path string) draws.Sprite {
	s := draws.NewSprite(path)
	s.SetScale(screenSizeX / s.W())
	s.SetPosition(screenSizeX/2, screenSizeY/2, draws.OriginCenterMiddle)
	return s
}
