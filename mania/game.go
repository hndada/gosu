package mania

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"log"
)

// scene: abstract; contents in the screen
// screen: mere image data after all; screen is the result

// 모든 scene에 sceneManager가 하는 일을 embed하면 없어도 되지 않을까?
// mania mode aka clavier
type Game struct {
	C      Chart
	Notes  []NoteImageInfo
	MaxTPS int

	scene       Scene
	scrollSpeed float64
	// input Input
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1600, 900
}

func (g *Game) Update(screen *ebiten.Image) error {
	if err := g.scene.Update(); err != nil {
		return err
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
}

type SceneMania struct {
	Notes        []NoteImageInfo
	C            Chart
	TickDuration float64
	scrollSpeed  float64
	Time         float64
	// Score float64
	// HP float64

	Cover *ebiten.Image
	// Audio
}

// 노트, 그냥 네모 그리고 색깔 채워넣기 하자
func (g *Game) NewSceneMania(c *Chart) (s *SceneMania) {
	s = &SceneMania{}
	const w = 70
	const h = 25
	s.C = *c
	s.Notes = make([]NoteImageInfo, len(c.Notes))
	s.scrollSpeed = g.scrollSpeed
	s.TickDuration = 1000 / float64(g.MaxTPS) // 스피드값 1 기준 초당 1000픽셀 내려오게 해야함
	for i, n := range c.Notes {
		x := float64(n.Key*w + 565)
		y := -float64(n.Time) * g.scrollSpeed
		s.Notes[i] = NoteImageInfo{x, y, w, h, noteColor(n, c.Keys)}
	}
	b, err := ioutil.ReadFile("./test/" + c.ImageFilename)
	if err != nil {
		panic(err)
	}
	cover, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	fmt.Println(cover.At(100, 100))
	s.Cover, _ = ebiten.NewImageFromImage(cover, ebiten.FilterDefault)
	// s.Cover.Fill(color.RGBA{139, 213, 238, 255})
	return
}

func (s *SceneMania) Update() error {
	s.Time += s.TickDuration
	for i := range s.Notes {
		s.Notes[i].y += s.TickDuration * s.scrollSpeed
	}
	return nil
}

// todo:범위 넘어간 애들은 Rect 안그리기
func (s *SceneMania) Draw(screen *ebiten.Image) {
	const ratio = float64(1600)/1920
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(ratio, ratio)
	op.ColorM.ChangeHSV(0, 1, 0.30)
	screen.DrawImage(s.Cover, op)

	ebitenutil.DrawRect(screen, 565, 0, 70*7, 900, color.RGBA{0, 0, 0, 180})
	ebitenutil.DrawRect(screen, 565, 730, 70*7, 10, color.RGBA{252, 106, 111, 255})
	for i := range s.Notes {
		ebitenutil.DrawRect(screen, s.Notes[i].x, s.Notes[i].y, s.Notes[i].w, s.Notes[i].h, s.Notes[i].clr)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTime: %.1fs", ebiten.CurrentFPS(), s.Time/1000))

}

type NoteImageInfo struct {
	x, y, w, h float64
	clr        color.RGBA
}

func noteColor(n Note, keys int) color.RGBA {
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

func main() {
	ebiten.SetWindowSize(1600, 900)
	ebiten.SetWindowTitle("gosu!")
	c := NewChart(`./test/test.osu`)
	g := &Game{}
	g.MaxTPS = 480
	g.scrollSpeed = 1.1
	ebiten.SetMaxTPS(g.MaxTPS)
	g.scene = g.NewSceneMania(c)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
