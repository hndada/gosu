package game

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/choose"
	"github.com/hndada/gosu/scene/play"
)

// All structs and variables in game package should be unexported
// since the game package is for being called at main via NewGame.
type game struct {
	fs.FS
	Scene
}
type Scene interface {
	Update() any
	Draw(screen draws.Image)
}

func NewGame(fsys fs.FS) *game {
	load(fsys)
	g := &game{
		FS:    fsys,
		Scene: nil,
	}
	return g
}
func (g *game) Update() (err error) {
	if g.Scene == nil {
		g.Scene = choose.NewScene()
		return
	}
	switch r := g.Scene.Update().(type) {
	case play.Return:
		g.Scene = choose.NewScene()
	case choose.Return:
		var scene Scene
		scene, err = play.NewScene(r.FS, r.Name, r.Mode, r.Mods, r.Replay)
		if err != nil {
			return
		}
		g.Scene = scene
	case error:
		err = r
	}
	return
}
func (g *game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(draws.Image{Image: screen})
}
func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scene.ScreenSizeX, scene.ScreenSizeY
}
