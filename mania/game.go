package mania

import (
	"bytes"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// Chart struct는 계산한 값 살리기 등을 위해서도 gob로 저장

// 곡선택: 맵정보패널
// 플레이
// 리플레이 실행 - 스코어/hp 시뮬레이터
type Game struct {
	GameState
	GameOptions
}

type GameState struct {
	Scene Scene
	input Input

	NextScene      Scene
	TransSceneFrom *ebiten.Image
	TransSceneTo   *ebiten.Image
	TransCountdown int
}

type Scene interface {
	Update(gs *GameState) error
	Draw(screen *ebiten.Image)
}
type Input interface{}

type GameOptions struct {
	MaxTPS      int
	ScrollSpeed float64
	KeysLayout  [][]ebiten.Key
	HitPosition float64
	// ScreenSize
}
type ImageInfo struct {
	x, y, w, h float64
	clr        color.RGBA
}

// scene: abstract; contents in the screen
// screen: mere image data after all; screen is the result

// 모든 scene에 sceneManager가 하는 일을 embed하면 없어도 되지 않을까?

// 곡선택 대비

// mp3, Scene에 저장; 연동
// sync with mp3, position
// input
// 점수계산: 1/n -> my score system

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1600, 900
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.TransCountdown == 0 {
		return g.Scene.Update(&g.GameState)
	}
	g.TransCountdown--
	if g.TransCountdown > 0 {
		return nil
	}
	// count down has just been from non-zero to zero
	g.Scene = g.NextScene
	g.NextScene = nil
	return nil
}

// scene의 Draw는 input으로 들어온 screen을 그리는 함수

func (g *Game) Draw(screen *ebiten.Image) {
	// _ = ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.CurrentFPS())) // 겹쳐버리는듯
	if g.TransCountdown == 0 {
		g.Scene.Draw(screen)
		return
	}
	var value float64
	var op ebiten.DrawImageOptions

	g.TransSceneFrom.Clear()
	g.Scene.Draw(g.TransSceneFrom)
	value = float64(g.TransCountdown) / 99 // todo: 변경 가능하게
	op = ebiten.DrawImageOptions{}
	// op.ColorM.Scale(1, 1, 1, alpha)
	op.ColorM.ChangeHSV(0, 1, value)
	_ = screen.DrawImage(g.TransSceneFrom, &op)

	g.TransSceneTo.Clear()
	// g.TransSceneTo.Fill(color.RGBA{128, 128, 0, 255}) // temp
	g.NextScene.Draw(g.TransSceneTo)
	value = 1 - float64(g.TransCountdown)/99 // todo: 변경 가능하게
	op = ebiten.DrawImageOptions{}
	// op.ColorM.Scale(1, 1, 1, alpha)
	op.ColorM.ChangeHSV(0, 1, value)
	_ = screen.DrawImage(g.TransSceneTo, &op)
}

type SceneMania struct { // aka Clavier
	Notes        []NoteImageInfo
	C            Chart
	TickDuration float64
	scrollSpeed  float64
	Time         float64
	// Score float64
	// HP float64

	Cover *ebiten.Image
	// Audio
	Done bool
}

// 노트, 그냥 네모 그리고 색깔 채워넣기
func NewSceneMania(op GameOptions, c *Chart) (s *SceneMania) {
	ebiten.SetWindowTitle(fmt.Sprintf("gosu - %s [%s]", c.Title, c.ChartName))
	s = &SceneMania{}
	const w = 70
	const noteHeight = 25
	s.C = *c
	s.Notes = make([]NoteImageInfo, len(c.Notes))
	s.scrollSpeed = op.ScrollSpeed
	s.TickDuration = 1000 / float64(op.MaxTPS) // 스피드값 1 기준 초당 1000픽셀 내려오게 해야함
	for i, n := range c.Notes {
		var y, h float64
		x := float64(n.Key*w + 565)
		switch n.Type {
		case TypeNote:
			y = -float64(n.Time)*op.ScrollSpeed + 1000
			h = noteHeight * op.ScrollSpeed
		case LNHead:
			y = -float64(n.Time2)*op.ScrollSpeed + 1000
			h = float64(n.Time2-n.Time+noteHeight) * op.ScrollSpeed
		}
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
	s.Cover, _ = ebiten.NewImageFromImage(cover, ebiten.FilterDefault)
	return
}

type NoteImageInfo = ImageInfo

func (s *SceneMania) Update(gs *GameState) error {
	var endTime float64 = 2.40 * 1000 // float64(s.C.Notes[len(s.C.Notes)-1].Time)
	s.Time += s.TickDuration
	for i := range s.Notes {
		// if s.Notes[i].y > 900 { continue }
		s.Notes[i].y += s.TickDuration * s.scrollSpeed
	}
	if s.Time > endTime {
		gs.NextScene = &SceneResult{}
		gs.TransCountdown = 99 // todo: set maxCount
	} // 끝났을 때
	return nil
}

type SceneResult struct {
	// score
	// hp graph
	// hit error deviation
}

func (s *SceneResult) Update(gs *GameState) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		panic("game over!")
	}
	// 키 입력 받으면 곡선택 scene으로 이동
	return nil
}

func (s *SceneResult) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Result")
}

// todo:범위 넘어간 애들은 Rect 안그리기 -> 오히려 fps 불안정
func (s *SceneMania) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	const ratio = float64(1600) / 1920
	op.GeoM.Scale(ratio, ratio)
	op.ColorM.ChangeHSV(0, 1, 0.30)

	screen.DrawImage(s.Cover, op)

	ebitenutil.DrawRect(screen, 565, 0, 70*7, 900, color.RGBA{0, 0, 0, 180})
	ebitenutil.DrawRect(screen, 565, 730, 70*7, 10, color.RGBA{252, 106, 111, 255})
	for i := range s.Notes {
		// if s.Notes[i].y>900 { continue }
		// if s.Notes[i].y+s.Notes[i].h<0 { continue }
		ebitenutil.DrawRect(screen, s.Notes[i].x, s.Notes[i].y, s.Notes[i].w, s.Notes[i].h, s.Notes[i].clr)
	}
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTime: %.1fs", ebiten.CurrentFPS(), s.Time/1000))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nTime: %.1fs", s.Time/1000))
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
	ebiten.SetWindowSize(1600, 900) // fixed in prototype
	ebiten.SetWindowTitle("gosu")
	c := NewChart(`./test/test_ln.osu`)
	g := &Game{}
	g.MaxTPS = 240
	g.ScrollSpeed = 1.33
	ebiten.SetMaxTPS(g.MaxTPS)
	ebiten.SetRunnableOnUnfocused(true)
	g.Scene = NewSceneMania(g.GameOptions, c)

	g.TransSceneFrom, _ = ebiten.NewImage(1600, 900, ebiten.FilterDefault)
	g.TransSceneTo, _ = ebiten.NewImage(1600, 900, ebiten.FilterDefault)

	f, err := os.Open("./test/" + c.AudioFilename)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
	<-done
}
