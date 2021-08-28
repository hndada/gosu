package gosu

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mania"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type ChartPanel struct {
	BodyText game.Sprite
	Body     game.LongSprite
	Left     game.Sprite
	Right    game.Sprite
	//idx   int
	Chart *mania.Chart
	BG    *ebiten.Image
	OpBG  *ebiten.DrawImageOptions
}

// Key - MusicName
// Level - ChartName
// X와 Y는 update에서 매번 새로 설정
func (s SceneSelect) NewChartPanel(c *mania.Chart) ChartPanel {
	const ChartPanelHeight = 40
	var cp ChartPanel
	cp.Body.Sprite = game.NewSprite(game.Skin.BoxMiddle)
	cp.Body.Vertical = false
	cp.Body.W = 450
	cp.Body.H = ChartPanelHeight
	{
		src := game.Skin.BoxLeft
		sprite := game.NewSprite(src)
		sprite.H = ChartPanelHeight
		scale := float64(sprite.H) / float64(src.Bounds().Dy())
		sprite.W = int(float64(src.Bounds().Dx()) * scale)
		cp.Left = sprite
	}
	{
		src := game.Skin.BoxRight
		sprite := game.NewSprite(src)
		sprite.H = ChartPanelHeight
		scale := float64(sprite.H) / float64(src.Bounds().Dy())
		sprite.W = int(float64(src.Bounds().Dx()) * scale)
		cp.Right = sprite
	}
	{
		rect := image.Rect(0, 0, 450, 40)
		img := image.NewRGBA(rect)
		PanelTextcolor := color.Black
		x, y := 20, 30
		point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(PanelTextcolor),
			Face: basicfont.Face7x13,
			Dot:  point,
		}
		t := fmt.Sprintf("(%dKey Lv %.1f) %s [%s]", c.KeyCount, c.Level, c.MusicName, c.ChartName)
		d.DrawString(t)

		src, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		sprite := game.NewSprite(src)
		sprite.W = 450
		sprite.H = 40
		cp.BodyText = sprite
	}
	cp.Chart = c
	//cp.BG, _ = c.Background()
	//cp.OpBG = game.BackgroundOp(s.ScreenSize, image.Pt(cp.BG.Size()))
	return cp
}

func (cp *ChartPanel) SetXY(x, y int) {
	cp.Body.X = x
	cp.BodyText.X = x
	cp.Left.X = x - cp.Left.W
	cp.Right.X = x + cp.Body.W

	cp.Body.Y = y
	cp.BodyText.Y = y
	cp.Left.Y = y
	cp.Right.Y = y
}

func (cp ChartPanel) Draw(screen *ebiten.Image) {
	cp.Left.Draw(screen)
	cp.Body.Draw(screen)
	cp.BodyText.Draw(screen)
	cp.Right.Draw(screen)
}

// type chartPanel struct {
// 	box   *ebiten.Image
// 	x, y  int // todo: sprite-ize
// 	op    *ebiten.DrawImageOptions
// 	chart *mania.Chart
// }

// func newChartPanel(c *mania.Chart) chartPanel {
// 	var cp chartPanel
// 	img := image.NewRGBA(image.Rect(0, 0, 450, 40))
// 	col := color.RGBA{200, 100, 0, 255}
// 	x, y := 20, 30
// 	point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
// 	d := &font.Drawer{
// 		Dst:  img,
// 		Src:  image.NewUniform(col),
// 		Face: basicfont.Face7x13,
// 		Dot:  point,
// 	}
// 	d.DrawString(fmt.Sprintf("(%dKey Lv %.1f) %s [%s]", c.KeyCount, c.Level, c.MusicName, c.ChartName))
// 	cp.box, _ = ebiten.NewImageFromImage(img, ebiten.FilterDefault)
// 	cp.op = &ebiten.DrawImageOptions{}
// 	cp.chart = c
// 	return cp
// }
