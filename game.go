package gosu

import (
	"image"
	"log"
	"os"
	"reflect"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mania"
	// _ "github.com/silbinarywolf/preferdiscretegpu"
)

var MaxTransCountDown int

// Game: path + Renderer
type Game struct {
	cwd            string // current working dir
	Scene          Scene
	NextScene      Scene
	TransSceneFrom *ebiten.Image
	TransSceneTo   *ebiten.Image
	TransCountdown int

	args game.TransSceneArgs
	// screenSize image.Point
}

type Scene interface {
	Ready() bool
	Update() error
	Draw(screen *ebiten.Image)           // Draws scene to screen
	Done(args *game.TransSceneArgs) bool // 모든 passed parameter는 Passed by Value.
}

func NewGame() *Game {
	const maxTPS = 60
	var err error

	g := &Game{}
	g.cwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	game.Settings.ScreenSize = image.Pt(800, 600)
	game.LoadSkin(g.cwd)
	mania.LoadSkin(g.cwd)

	ebiten.SetWindowSize(game.Settings.ScreenSize.X, game.Settings.ScreenSize.Y)
	g.Scene = newSceneSelect(g.cwd, game.Settings.ScreenSize)

	g.args = game.TransSceneArgs{}
	ebiten.SetWindowTitle("gosu")
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetMaxTPS(maxTPS)

	g.TransSceneFrom, _ = ebiten.NewImage(game.Settings.ScreenSize.X, game.Settings.ScreenSize.Y, ebiten.FilterDefault)
	g.TransSceneTo, _ = ebiten.NewImage(game.Settings.ScreenSize.X, game.Settings.ScreenSize.Y, ebiten.FilterDefault)
	MaxTransCountDown = ebiten.MaxTPS() * 4 / 5
	return g
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.TransCountdown <= 0 { // == 0
		if g.Scene.Done(&g.args) {
			switch g.Scene.(type) {
			case *sceneSelect:
				switch g.args.Next {
				case "mania.Scene":
					v := reflect.ValueOf(g.args.Args)
					chart := v.FieldByName("Chart").Interface().(*mania.Chart)
					mods := v.FieldByName("Mods").Interface().(mania.Mods)
					p := v.FieldByName("ScreenSize").Interface().(image.Point)
					s2 := mania.NewScene(chart, mods, p, g.cwd)
					g.ChangeScene(s2)
				}
			case *mania.Scene:
				s2 := newSceneSelect(g.cwd, game.Settings.ScreenSize) // temp: 매번 새로 만들 필요는 없음
				g.ChangeScene(s2)
			default:
				panic("not reach")
			}
			g.args = game.TransSceneArgs{}
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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.TransCountdown == 0 {
		if g.Scene.Ready() {
			g.Scene.Draw(screen)
		}
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
	return game.Settings.ScreenSize.X, game.Settings.ScreenSize.Y
}

func (g *Game) ChangeScene(s Scene) {
	g.NextScene = s
	g.TransCountdown = MaxTransCountDown
}
