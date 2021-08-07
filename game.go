package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mania"

	// _ "github.com/silbinarywolf/preferdiscretegpu"
	"path/filepath"
)

var MaxTransCountDown int

const gosuPath = `E:\gosu\`

// Game: path + Renderer
type Game struct {
	path           string
	Scene          Scene
	NextScene      Scene
	TransSceneFrom *ebiten.Image
	TransSceneTo   *ebiten.Image
	TransCountdown int
}

// Scene: an actual thing that control the game
type Scene interface {
	Init()
	Update() error
	Draw(screen *ebiten.Image) // Draws scene to screen
	Done() bool
}

func NewGame() *Game {
	g := &Game{}
	g.path = gosuPath

	game.LoadSettings()
	mania.ResetSettings()
	mania.LoadSpriteMap(filepath.Join(g.path, "Skin"))
	g.Scene = g.NewSceneSelect()

	p := game.ScreenSize()
	g.TransSceneFrom, _ = ebiten.NewImage(p.X, p.Y, ebiten.FilterDefault)
	g.TransSceneTo, _ = ebiten.NewImage(p.X, p.Y, ebiten.FilterDefault)
	MaxTransCountDown = game.MaxTPS() * 4 / 5

	ebiten.SetWindowTitle("gosu")
	ebiten.SetRunnableOnUnfocused(true)
	return g
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.TransCountdown <= 0 { // == 0
		if g.Scene.Done() {
			g.ChangeScene(g.NewSceneSelect()) // temp
		}
		return g.Scene.Update()
	}
	g.TransCountdown--
	if g.TransCountdown > 0 {
		return nil
	}
	// count down has just been from non-zero to zero
	g.Scene = g.NextScene
	g.NextScene = nil
	g.Scene.Init()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.TransCountdown == 0 {
		g.Scene.Draw(screen)
		return
	}
	var value float64
	{
		value = float64(g.TransCountdown) / float64(MaxTransCountDown)
		g.TransSceneFrom.Clear()
		g.Scene.Draw(g.TransSceneFrom)
		op := ebiten.DrawImageOptions{}
		op.ColorM.ChangeHSV(0, 1, value)
		screen.DrawImage(g.TransSceneFrom, &op)
	}
	{
		value = 1 - value
		g.TransSceneTo.Clear()
		g.NextScene.Draw(g.TransSceneTo)
		op := ebiten.DrawImageOptions{}
		op.ColorM.ChangeHSV(0, 1, value)
		screen.DrawImage(g.TransSceneTo, &op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return game.ScreenSize().X, game.ScreenSize().Y
}

func (g *Game) ChangeScene(s Scene) {
	g.NextScene = s
	g.TransCountdown = MaxTransCountDown
}

// BasePlayScene (base struct), PlayScene (interface)
// PlayScene, 각 mode 패키지에다가 구현해야 할까?
