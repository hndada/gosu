package gosu

import (
	"os"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/engine/scene"
	"github.com/hndada/gosu/mania"
	// _ "github.com/silbinarywolf/preferdiscretegpu"
)

var cwd string            // current working dir
var charts []*mania.Chart // temp

// background goes lazy loaded
type Game struct {
	scene   scene.Scene
	args    *scene.Args
	changer *scene.Changer
}

func NewGame() *Game {
	ebiten.SetWindowSize(common.Settings.ScreenSize.X, common.Settings.ScreenSize.Y)
	ebiten.SetWindowTitle("gosu")
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetMaxTPS(common.Settings.MaxTPS)

	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	common.LoadSkin(cwd)
	mania.LoadSkin(cwd)
	charts = loadCharts(cwd)
	sceneSelect = newSceneSelect(cwd)

	g := &Game{}
	g.scene = sceneSelect
	g.args = &scene.Args{}
	g.changer = scene.NewChanger()
	return g
}

// whether changer or scene goes updated
func (g *Game) Update() error {
	if !g.changer.Done() {
		return g.changer.Update()
	}
	if g.scene.Close(g.args) {
		switch g.scene.(type) {
		case *SceneSelect:
			switch g.args.Next {
			case "mania.Scene":
				v := reflect.ValueOf(g.args.Args)
				chart := v.FieldByName("Chart").Interface().(*mania.Chart)
				mods := v.FieldByName("Mods").Interface().(mania.Mods)
				s2 := mania.NewScene(chart, mods, cwd)
				g.changer.Change(g.scene, s2)
				g.scene = s2
			}
		case *mania.Scene:
			updateCharts(cwd)
			sceneSelect.close = false
			g.changer.Change(g.scene, sceneSelect)
			g.scene = sceneSelect
		default:
			panic("not reach")
		}
		g.args = &scene.Args{} // TODO: must?
	} else {
		return g.scene.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.changer.Done() {
		g.scene.Draw(screen)
	} else {
		g.changer.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return common.Settings.ScreenSize.X, common.Settings.ScreenSize.Y
}
