package mania

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
)

// scene: abstract; contents in the screen
// screen: mere image data after all

// screen is the result

const (
	ScreenWidth  = 320
	ScreenHeight = 240
)

func init() {
	ebiten.SetMaxTPS(240)
}

// 모든 scene에 sceneManager가 하는 일을 embed하면 없어도 되지 않을까?
type Game struct {
	// speed float64
	scene ScenePlay
	// scene Scene
	// input Input
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update(screen *ebiten.Image) error {
	fmt.Println("update")
	if err := g.scene.Update(screen); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	fmt.Println("draw")
	g.scene.Draw(screen)
}

type ScenePlay struct {
	Notes []NoteImageInfo
	C     Chart
	// Notes []ebiten.Image // for showing
	// Score float64
	// HP float64
}

// 노트, 그냥 네모 그리고 색깔 채워넣기 하자
func (g *Game) NewScenePlay(c *Chart) (s ScenePlay) {
	const w = 30
	const h = 10
	s.C = *c
	s.Notes = make([]NoteImageInfo, len(c.Notes))
	for i, n := range c.Notes {
		x := float64(n.Key*w + 200)
		y := -float64(n.Time) + 400 // g.speed
		// ebitenutil.DrawRect(&s.Notes[i], x, y, w, h, noteColor(n))
		s.Notes[i] = NoteImageInfo{x, y, w, h, noteColor(n)}
		// s.Notes[i].DrawImage()
	}
	fmt.Println(s.Notes)
	return
}

type NoteImageInfo struct {
	x, y, w, h float64
	clr        color.RGBA
}

func noteColor(n Note) color.RGBA {
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

func (s *ScenePlay) Update(screen *ebiten.Image) error {
	// const scrollSpeed = 10
	// unitTime := 1 / ebiten.CurrentTPS()
	// unitMove := unitTime * scrollSpeed // distance
	// for i := range s.Notes {
	// 	s.Notes[i].y += unitMove
	// }
	return nil
}

func (s *ScenePlay) Draw(screen *ebiten.Image) {
	// var err error
	for i := range s.Notes {
		// if err = s.Notes[i].DrawImage(screen, &ebiten.DrawImageOptions{}); err != nil {
		// 	panic(err)
		// }
		ebitenutil.DrawRect(screen, s.Notes[i].x, s.Notes[i].y, s.Notes[i].w, s.Notes[i].h, s.Notes[i].clr)
	}
}

// type Scene interface {
// 	Update(screen *ebiten.Image) error
// 	Draw(screen *ebiten.Image)
// }
