package box_test

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/hndada/gosu/ebitenui"
	"image"
	"image/color"
	"log"
	"testing"
)

type Game struct {
	tbox        *ebiten.Image
	cbox        *ebitenui.Checkbox
	button      *ebitenui.Button
	buttonCount int
}

func (g *Game) Update(screen *ebiten.Image) error {
	g.cbox.Update()
	g.button.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{220, 80, 30, 255})
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(30, 30)
		screen.DrawImage(g.tbox, op)
	}
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(30, 65)
		screen.DrawImage(ebitenui.ImgCheckedBox, op)
		op.GeoM.Translate(40, 0)
		screen.DrawImage(ebitenui.ImgUncheckedBox, op)
	}
	g.cbox.Draw(screen)
	g.button.Draw(screen)
	text.Draw(screen, fmt.Sprintf("button hit: %d", g.buttonCount), ebitenui.MplusNormalFont, 30, 150, color.Black)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func TestMain(m *testing.M) { // should be always main
	g := &Game{}
	ebiten.SetWindowSize(320, 240)

	t := ebitenui.RenderText("text box", ebitenui.MplusNormalFont, color.Black)
	g.tbox = ebitenui.RenderTextBox(t, color.White)
	g.cbox = ebitenui.NewCheckbox("check box", image.Point{150, 65})

	bt := ebitenui.RenderText("Button", ebitenui.MplusNormalFont, color.Opaque)
	g.button = &ebitenui.Button{
		MinPt: image.Point{30, 95},
		Image: ebitenui.RenderTextBox(bt, color.RGBA{192, 18, 112, 255}),
	}
	g.button.SetOnPressed(func(b *ebitenui.Button) {
		g.buttonCount++
	})

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}