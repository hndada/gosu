package mode

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

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
	DefaultBackground draws.Sprite
	CursorSprites     [2]draws.Sprite // 0: cursor // 1: additive cursor
	// CursorTailSprite   Sprite
	ChartInfoBoxSprite draws.Sprite // Todo: various box sprite
)

func LoadGeneralSkin() {
	DefaultBackground = draws.Sprite{
		I:      draws.NewImage("skin/default-bg.jpg"),
		Filter: ebiten.FilterLinear,
	}
	DefaultBackground.SetWidth(screenSizeX)
	DefaultBackground.SetCenterY(screenSizeY / 2)

	for i, name := range []string{"menu-cursor.png", "menu-cursor-additive.png"} {
		s := draws.Sprite{
			I:      draws.NewImage(fmt.Sprintf("skin/cursor/%s", name)),
			Filter: ebiten.FilterLinear,
		}
		s.ApplyScale(CursorScale)
		CursorSprites[i] = s
	}

	purple := color.RGBA{172, 49, 174, 255}
	white := color.RGBA{255, 255, 255, 128}
	const border = 3
	w := int(ChartInfoBoxWidth)
	h := int(ChartInfoBoxHeight)

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), &image.Uniform{purple}, image.Point{}, draw.Src)
	inRect := image.Rect(border, border, w-border, h-border)
	draw.Draw(img, inRect, &image.Uniform{white}, image.Point{}, draw.Src)
	ChartInfoBoxSprite = draws.Sprite{
		I: ebiten.NewImageFromImage(img),
		W: float64(w),
		H: float64(h),
		X: screenSizeX - float64(w) + chartInfoBoxshrink,
		// Y is not fixed.
	}
}
