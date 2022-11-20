package gosu

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
)

// All structs and variables in game package should be unexported
// since the game package is for being called at main via NewGame.
type game struct {
	fs.FS
	scene.Scene
}

//	func NewGame(fsys fs.FS) *game {
//		load(fsys)
//		g := &game{
//			FS:    fsys,
//			Scene: choose.NewScene(),
//		}
//		return g
//	}
func load(fsys fs.FS) {
	settings, err := fs.ReadFile(fsys, "settings.toml")
	if err != nil {
		fmt.Println("No custom setting file detected.")
	}
	mode.UserSettings.Load(string(settings))
	piano.UserSettings.Load(string(settings))
	drum.UserSettings.Load(string(settings))
	// scene.UserSettings.Load(string(settings))

	skinFS, err := fs.Sub(fsys, "skin")
	if err != nil {
		fmt.Println(err)
	}
	mode.UserSkin.Load(skinFS)
	piano.UserSkins.Load(skinFS)
	drum.UserSkin.Load(skinFS)
	// scene.UserSkin.Load(skinFS)
}
func (g *game) Update() (err error) {
	// switch r := g.Scene.Update().(type) {
	// case error:
	// 	err = r
	// case choose.Return:
	// 	var scene scene.Scene
	// 	scene, err = play.NewScene(r.FS, r.Name, r.Mode, r.Mods, r.Replay)
	// 	if err != nil {
	// 		return
	// 	}
	// 	g.Scene = scene
	// case mode.Result:
	// 	g.Scene = choose.NewScene()
	// }
	g.Scene.Update()
	return
}
func (g *game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(draws.Image{Image: screen})
}
func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return mode.ScreenSizeX, mode.ScreenSizeY
}
