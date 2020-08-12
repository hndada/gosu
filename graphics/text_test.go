package graphics

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"image/color"
	"log"
	"testing"
)

type Game struct {
	textBox *ebiten.Image
	op      ebiten.DrawImageOptions
}

func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	// ebitenutil.DebugPrintAt(screen, "a", 100, 100) // covered
	screen.Fill(color.RGBA{220, 80, 30, 255})
	ebitenutil.DebugPrint(screen, "hey")
	screen.DrawImage(g.textBox, &g.op)
	text.Draw(screen, "Temp", FontVarelaNormal, 100, 100, color.Black)

	// b := text.BoundString(FontVarelaNormal, "heh")
	// img, _ := ebiten.NewImage(b.Dx(), b.Dy(), ebiten.FilterDefault)
	// text.Draw(img, "heh", FontVarelaNormal, 0, 200, color.Black)
	// screen.DrawImage(img, &ebiten.DrawImageOptions{})
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

// 언제나 main이어야 한다
func TestMain(m *testing.M) {
	g := &Game{}
	ebiten.SetWindowSize(320, 240)
	text := DrawText("hello", FontVarelaNormal, color.Black)
	g.textBox = DrawTextBox(text, color.White)
	g.op.GeoM.Translate(200, 150)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
