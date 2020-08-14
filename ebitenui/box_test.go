package ebitenui

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"image"
	"image/color"
	"log"
	"testing"
)

type Game struct {
	tbox        *ebiten.Image
	cbox        *Checkbox
	button      *Button
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
		screen.DrawImage(imgCheckedBox, op)
		op.GeoM.Translate(40, 0)
		screen.DrawImage(imgUncheckedBox, op)
	}
	g.cbox.Draw(screen)
	g.button.Draw(screen)
	text.Draw(screen, fmt.Sprintf("button hit: %d", g.buttonCount), MplusNormalFont, 30, 150, color.Black)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func TestMain(m *testing.M) { // should be always main
	g := &Game{}
	ebiten.SetWindowSize(320, 240)

	t := RenderText("text box", MplusNormalFont, color.Black)
	g.tbox = RenderTextBox(t, color.White)
	g.cbox = NewCheckbox("check box", image.Point{150, 65})

	bt := RenderText("Button", MplusNormalFont, color.Opaque)
	g.button = &Button{
		MinPt: image.Point{30, 95},
		Image: RenderTextBox(bt, color.RGBA{192, 18, 112, 255}),
	}
	g.button.SetOnPressed(func(b *Button) {
		g.buttonCount++
	})

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}