package graphics

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
	text.Draw(screen, fmt.Sprintf("button hit: %d", g.buttonCount), mplusNormalFont, 30, 150, color.Black)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

// 언제나 main이어야 한다
func TestMain(m *testing.M) {
	g := &Game{}
	ebiten.SetWindowSize(320, 240)

	t := DrawText("text box", mplusNormalFont, color.Black)
	g.tbox = DrawTextBox(t, color.White)
	g.cbox = NewCheckbox("check box", image.Point{150, 65})

	bt := DrawText("Button", mplusNormalFont, color.Opaque)
	g.button = &Button{
		MinPt: image.Point{30, 95},
		Image: DrawTextBox(bt, color.RGBA{192, 18, 112, 255}),
	}
	g.button.SetOnPressed(func(b *Button) {
		g.buttonCount++
	})

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func f64Pt(p image.Point) (float64, float64) {
	return float64(p.X), float64(p.Y)
}
