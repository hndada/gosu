package mania

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
)

// scene: abstract; contents in the screen
// screen: mere image data after all; screen is the result

// 모든 scene에 sceneManager가 하는 일을 embed하면 없어도 되지 않을까?
type Game struct {
	C      Chart
	Notes  []NoteImageInfo
	MaxTPS int

	// scrollSpeed float64
	scene ScenePlay
	// input Input
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

func (g *Game) Update(screen *ebiten.Image) error {
	if err := g.scene.Update(screen); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (s *ScenePlay) Update(screen *ebiten.Image) error {
	// const scrollSpeed = 0.01
	// unitTime := 1 / ebiten.CurrentTPS()
	// unitMove := unitTime * scrollSpeed // distance
	// fmt.Println(ebiten.CurrentTPS())
	f := float64(ebiten.MaxTPS())
	for i := range s.Notes {
		s.Notes[i].y += 240 / f
	}
	return nil
}

// 범위 넘어간 애들은 Rect 안그리기
func (s *ScenePlay) Draw(screen *ebiten.Image) {
	// fmt.Println("draw")
	for i := range s.Notes {
		// if err = s.Notes[i].DrawImage(screen, &ebiten.DrawImageOptions{}); err != nil {
		// 	panic(err)
		// }
		ebitenutil.DrawRect(screen, s.Notes[i].x, s.Notes[i].y, s.Notes[i].w, s.Notes[i].h, s.Notes[i].clr)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprint(ebiten.CurrentFPS()))
}

func main() {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("gosu!")
	c := NewChart(`C:\Users\hndada\Documents\GitHub\hndada\gosu\mania\test.osu`)
	g := &Game{}
	g.MaxTPS = 480
	g.scene=NewScenePlay(c)
	ebiten.SetMaxTPS(g.MaxTPS)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// 노트, 그냥 네모 그리고 색깔 채워넣기 하자
func NewScenePlay(c *Chart) (s ScenePlay) {
	const w = 48
	const h = 10
	s.C = *c
	s.Notes = make([]NoteImageInfo, len(c.Notes))
	for i, n := range c.Notes {
		x := float64(n.Key*w + 200)
		y := -float64(n.Time)/3 + 400
		s.Notes[i] = NoteImageInfo{x, y, w, h, noteColor(n)}
	}
	return
}

type NoteImageInfo struct {
	x, y, w, h float64
	clr        color.RGBA
}

func noteColor(n Note) color.RGBA {
	switch n.Key {
	case 0, 2, 4, 6:
		return color.RGBA{239, 243, 247, 0xff} // white
	case 1, 5:
		return color.RGBA{66, 211, 247, 0xff} // blue
	case 3:
		return color.RGBA{255, 203, 82, 0xff} // yellow
	}
	panic("not reach")
}

type ScenePlay struct {
	Notes []NoteImageInfo
	C     Chart
	// Notes []ebiten.Image // for showing
	// Score float64
	// HP float64
}

type Scene interface {
	Update(screen *ebiten.Image) error
	Draw(screen *ebiten.Image)
}



//
// func (s *ScenePlay) Update(screen *ebiten.Image) error {
// 	return nil
// }
//
// func (s *ScenePlay) Draw(screen *ebiten.Image) {
//
// }
