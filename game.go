package gosu

import (
	"reflect"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mania"

	// _ "github.com/silbinarywolf/preferdiscretegpu"
	"path/filepath"
)

var MaxTransCountDown int

const gosuPath = `E:\gosu\`

// BasePlayScene (base struct), PlayScene (interface)
// PlayScene, 각 mode 패키지에다가 구현해야 할까?
// Game: path + Renderer
type Game struct {
	cwd            string // current working dir
	path           string
	Scene          Scene
	NextScene      Scene
	TransSceneFrom *ebiten.Image
	TransSceneTo   *ebiten.Image
	TransCountdown int

	args game.TransSceneArgs
}

// Scene: an actual thing that control the game
// 다음 scene에게 필요한 정보를 넘겨줘야 함
type Scene interface {
	Init()
	Update() error
	Draw(screen *ebiten.Image)           // Draws scene to screen
	Done(args *game.TransSceneArgs) bool // 모든 passed parameter는 Passed by Value.
}

func NewGame() *Game {
	g := &Game{}
	g.path = gosuPath

	game.LoadSettings()
	mania.ResetSettings()
	mania.LoadSpriteMap(filepath.Join(g.path, "Skin"))
	g.Scene = newSceneSelect(g.path)

	p := game.ScreenSize()
	g.TransSceneFrom, _ = ebiten.NewImage(p.X, p.Y, ebiten.FilterDefault)
	g.TransSceneTo, _ = ebiten.NewImage(p.X, p.Y, ebiten.FilterDefault)
	MaxTransCountDown = game.MaxTPS() * 4 / 5

	g.args = game.TransSceneArgs{}
	ebiten.SetWindowTitle("gosu")
	ebiten.SetRunnableOnUnfocused(true)
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
					s2 := mania.NewScene(chart, mods)
					g.ChangeScene(s2)
				}
			case *mania.Scene:
				s2 := newSceneSelect(g.path) // temp: 매번 새로 만들 필요는 없음
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

func (g Game) CWD() string {
	return g.cwd
}
