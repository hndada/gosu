package game

import (
	"github.com/hajimehoshi/ebiten"
)

// mp3 플레이어, Scene에 저장; 연동
// sync with mp3, position
// 곡선택: 맵정보패널

// 플레이: input (ebiten으로 간단히, 나중에 별도 라이브러리.)
// 점수계산: 1/n -> my score system
// 리플레이 실행 - 스코어/hp 시뮬레이터
type Game struct {
	State
	Options
}

type State struct {
	Scene          Scene
	NextScene      Scene
	TransSceneFrom *ebiten.Image
	TransSceneTo   *ebiten.Image
	// TransCountdown int
	TransLifetime float64
	Loading       bool

	input Input
}

type Options struct {
	ScrollSpeed  float64
	KeysLayout   map[int][]ebiten.Key
	MaxTPS       int
	ScreenWidth  int
	ScreenHeight int
	HitPosition  float64 // object which is now set at 'options'

	DimValue  int
	VolumeSFX int
	VolumeBGM int
	Skin      *Skin
}

// scene이 game을 control하는 주체
type Scene interface {
	Update(g *Game) error
	Draw(screen *ebiten.Image)
}

func (g *Game) Tick() float64             { return 1000 / float64(g.MaxTPS) }
func (g *Game) MaxTransLifetime() float64 { return 600 } // 모든 time 관련 단위는 ms
func (g *Game) Update(screen *ebiten.Image) error {
	if g.Loading { return nil }
	if int(g.TransLifetime) == 0 {
		return g.Scene.Update(g)
	}
	g.TransLifetime -= 4.2 // g.Tick()
	// g.TransLifetime--
	if g.TransLifetime > 0 {
		return nil
	}
	// count down has just been from non-zero to zero
	g.Scene = g.NextScene
	g.NextScene = nil
	return nil
}

// scene의 Draw는 input으로 들어온 screen을 그리는 함수
func (g *Game) Draw(screen *ebiten.Image) {
	//  ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.CurrentFPS())) // 겹쳐버리는듯
	if g.Loading {
		return
	}
	if int(g.TransLifetime) == 0 {
		g.Scene.Draw(screen)
		return
	}
	var value float64
	var op ebiten.DrawImageOptions

	g.TransSceneFrom.Clear()
	g.Scene.Draw(g.TransSceneFrom)
	value = g.TransLifetime / g.MaxTransLifetime()
	op = ebiten.DrawImageOptions{}
	// op.ColorM.Scale(1, 1, 1, alpha)
	op.ColorM.ChangeHSV(0, 1, value)
	screen.DrawImage(g.TransSceneFrom, &op)

	g.TransSceneTo.Clear()
	g.NextScene.Draw(g.TransSceneTo)
	value = 1 - value
	op = ebiten.DrawImageOptions{}
	// op.ColorM.Scale(1, 1, 1, alpha)
	op.ColorM.ChangeHSV(0, 1, value)
	screen.DrawImage(g.TransSceneTo, &op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenWidth, g.ScreenHeight
}

func NewGame() (g *Game) {
	g = &Game{}
	g.State = State{
		Scene: &SceneTitle{},
	}
	g.Options = Options{ // todo: load settings
		MaxTPS:      240,
		ScrollSpeed: 1.33,
		HitPosition: 730,
	}
	g.ScreenWidth = 1600
	g.ScreenHeight = 900
	g.TransSceneFrom, _ = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight, ebiten.FilterDefault)
	g.TransSceneTo, _ = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight, ebiten.FilterDefault)
	return
}
