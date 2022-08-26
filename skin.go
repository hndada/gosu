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
	ScoreSprites       [10]draws.Sprite
	SignSprites        [3]draws.Sprite
)

const (
	CursorSpriteBase = iota
	CursorSpriteAdditive
	CursorSpriteTrail
)

func LoadGeneralSkin() {
	{
		s := draws.NewSprite("skin/default-bg.jpg")
		scale := screenSizeX / s.W()
		s.SetScale(scale, scale, ebiten.FilterLinear)
		s.SetPosition(screenSizeX/2, screenSizeY/2, draws.OriginCenter)
		DefaultBackground = s
	}
	// Todo: cursor may have outer circle
	names := []string{"menu-cursor", "menu-cursor-additive", "cursortrail"}
	for i, name := range names {
		s := draws.NewSprite(fmt.Sprintf("skin/cursor/%s.png", name))
		s.SetScale(CursorScale, CursorScale, ebiten.FilterLinear)
		s.SetPosition(screenSizeX/2, screenSizeY/2, draws.OriginCenter)
		CursorSprites[i] = s
	}
	{
		s := draws.NewSprite("skin/box.png")
		scaleW := ChartInfoBoxWidth / s.W()
		scaleH := ChartInfoBoxHeight / s.H()
		s.SetScale(scaleW, scaleH, ebiten.FilterLinear)
		offset := -(ChartInfoBoxWidth - chartInfoBoxshrink)
		s.SetPosition(screenSizeX+offset, screenSizeY/2, draws.OriginRightCenter)
		DefaultBackground = s
	}
	// purple := color.RGBA{172, 49, 174, 255}
	// white := color.RGBA{255, 255, 255, 128}
	// const border = 3
	// w := int(ChartInfoBoxWidth)
	// h := int(ChartInfoBoxHeight)

	// img := image.NewRGBA(image.Rect(0, 0, w, h))
	// draw.Draw(img, img.Bounds(), &image.Uniform{purple}, image.Point{}, draw.Src)
	// inRect := image.Rect(border, border, w-border, h-border)
	// draw.Draw(img, inRect, &image.Uniform{white}, image.Point{}, draw.Src)
	// ChartInfoBoxSprite = draws.Sprite{
	// 	I: ebiten.NewImageFromImage(img),
	// 	W: float64(w),
	// 	H: float64(h),
	// 	X: screenSizeX - float64(w) + chartInfoBoxshrink,
	// 	// Y is not fixed.
	// }
	for i := 0; i < 10; i++ {
		s := draws.NewSprite(fmt.Sprintf("skin/score/%d.png", i))
		s.SetScale(ScoreScale, ScoreScale, ebiten.FilterLinear)
		s.SetPosition(screenSizeX, 0, draws.OriginRightTop)
		ScoreSprites[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.NewSprite(fmt.Sprintf("skin/score/%s.png", name))
		s.SetScale(ScoreScale, ScoreScale, ebiten.FilterLinear)
		s.SetPosition(screenSizeX, 0, draws.OriginRightTop)
		SignSprites[i] = s
	}
}
