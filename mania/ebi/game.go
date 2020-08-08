package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/mania"
	"image/color"
	"log"
)

func init() {
	ebiten.SetMaxTPS(240)
}

type Game struct {
	C     mania.Chart
	Notes []NoteImageInfo
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

func (g *Game) Update(screen *ebiten.Image) error {
	// const scrollSpeed = 0.01
	// unitTime := 1 / ebiten.CurrentTPS()
	// unitMove := unitTime * scrollSpeed // distance
	for i := range g.Notes {
		g.Notes[i].y += 1
	}
	return nil
}

// 범위 넘어간 애들은 Rect 안그리기
func (g *Game) Draw(screen *ebiten.Image) {
	// fmt.Println("draw")
	for i := range g.Notes {
		// if err = g.Notes[i].DrawImage(screen, &ebiten.DrawImageOptions{}); err != nil {
		// 	panic(err)
		// }
		ebitenutil.DrawRect(screen, g.Notes[i].x, g.Notes[i].y, g.Notes[i].w, g.Notes[i].h, g.Notes[i].clr)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprint(ebiten.CurrentFPS()))
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("gosu!")
	c := mania.NewChart(`C:\Users\hndada\Documents\GitHub\hndada\gosu\mania\test.osu`)
	g := NewGame(c)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

func NewGame(c *mania.Chart) (g Game) {
	const w = 48
	const h = 10
	g.C = *c
	g.Notes = make([]NoteImageInfo, len(c.Notes))
	for i, n := range c.Notes {
		x := float64(n.Key*w + 200)
		y := -float64(n.Time)/3 + 400
		g.Notes[i] = NoteImageInfo{x, y, w, h, noteColor(n)}
	}
	return
}

//
type NoteImageInfo struct {
	x, y, w, h float64
	clr        color.RGBA
}

func noteColor(n mania.Note) color.RGBA {
	switch n.Key {
	case 0, 2, 4, 6: // white
		return color.RGBA{239, 243, 247, 0xff}
	case 1, 5: // blue
		return color.RGBA{66, 211, 247, 0xff}
	case 3: // yellow
		return color.RGBA{255, 203, 82, 0xff}
	}
	panic("not reach")
}
