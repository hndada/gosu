package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
)

type Game struct {
	Scene
}
type Scene interface {
	Update() any
	Draw(screen *ebiten.Image)
}

var sceneSelect *SceneSelect

func NewGame() *Game {
	piano.LoadSkin()
	sceneSelect = NewSceneSelect()
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(WindowSizeX, WindowSizeY)
	ebiten.SetMaxTPS(mode.MaxTPS)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	g := &Game{
		Scene: sceneSelect,
	}
	return g
}

type SelectToPlayArgs struct {
	Path   string
	Mode   int
	Replay *osr.Format
	Play   bool
}

func (g *Game) Update() error {
	args := g.Scene.Update()
	if args == nil {
		return nil
	}
	switch args := args.(type) {
	case mode.PlayToResultArgs:
		// Todo: selectResult
		g.Scene = sceneSelect
	case SelectToPlayArgs:
		switch args.Mode {
		// case args.Mode&mode.ModePiano != 0:
		case mode.ModePiano4, mode.ModePiano7:
			var err error
			g.Scene, err = piano.NewScenePlay(args.Path, args.Replay, args.Play)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}

//	func a(args *Args) {
//		args2 := reflect.ValueOf(args)
//
// from := args2.FieldByName("From").String()
//
//		NewSceneResult()
//		args2.FieldByName("Result")
//	}
