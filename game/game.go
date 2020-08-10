package game

import (
	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	State
	Options
}

type State struct {
	Scene          Scene
	NextScene      Scene
	TransSceneFrom *ebiten.Image
	TransSceneTo   *ebiten.Image
	TransCountdown int

	input Input
}

type Options struct {
	ScrollSpeed  float64
	KeysLayout   map[int][]ebiten.Key
	MaxTPS       int
	ScreenWidth  int
	ScreenHeight int
	HitPosition  float64 // skin -> option
}

type Scene interface {
	Update(gs *State) error
	Draw(screen *ebiten.Image)
}
type Input interface{}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenWidth, g.ScreenHeight
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.TransCountdown == 0 {
		return g.Scene.Update(&g.State)
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
	//  ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.CurrentFPS())) // 겹쳐버리는듯
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
	screen.DrawImage(g.TransSceneFrom, &op)

	g.TransSceneTo.Clear()
	// g.TransSceneTo.Fill(color.RGBA{128, 128, 0, 255}) // temp
	g.NextScene.Draw(g.TransSceneTo)
	value = 1 - float64(g.TransCountdown)/99 // todo: 변경 가능하게
	op = ebiten.DrawImageOptions{}
	// op.ColorM.Scale(1, 1, 1, alpha)
	op.ColorM.ChangeHSV(0, 1, value)
	screen.DrawImage(g.TransSceneTo, &op)
}

func NewGame() (g *Game) {
	g = &Game{}
	// g.State = State{
	// 	Scene:,
	// }
	g.Options = Options{
		MaxTPS:      240,
		ScrollSpeed: 1.33,
		// KeysLayout
		HitPosition: 730,
	}
	g.ScreenWidth = 1600
	g.ScreenHeight = 900
	g.TransSceneFrom, _ = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight, ebiten.FilterDefault)
	g.TransSceneTo, _ = ebiten.NewImage(g.ScreenWidth, g.ScreenHeight, ebiten.FilterDefault)
	return
}

// type ImageInfo struct {
// 	x, y, w, h float64
// 	clr        color.RGBA
// }
