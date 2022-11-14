package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	field *ebiten.Image
}

var tick int

func main() {
	g := Game{}
	s := ebiten.NewImage(50, 50)
	// s.Fill(color.NRGBA{128, 50, 30, 255})
	// s.Fill(color.NRGBA{195, 234, 131, 255})
	s.Fill(color.NRGBA{128, 0, 128, 255})
	g.field = s
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
func (g Game) Update() error {
	tick++
	return nil
}
func (g Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{128, 128, 128, 255})
	op := &ebiten.DrawImageOptions{}
	// op.ColorM.RotateHue(float64(tick) / 60)
	switch {
	case tick%240 < 120:
		op.ColorM.ChangeHSV(0, 0, 1)
	default:
		op.ColorM.ChangeHSV(0, 1, 1)
	}
	screen.DrawImage(g.field, op)
	ebitenutil.DebugPrint(
		screen, fmt.Sprintf("%4.2f\n%4.2f",
			ebiten.ActualFPS(),
			ebiten.ActualTPS()))
}

const (
	screenSizeX = 320
	screenSizeY = 240
)

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}
