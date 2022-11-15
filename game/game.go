package game

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
)

// All structs and variables in game package should be unexported
// since the game package is for being called at main via NewGame.
type game struct {
	fs.FS
	scene.Scene
}

func NewGame(fsys fs.FS) *game {
	scene.Load(fsys)
	g := &game{
		FS:    fsys,
		Scene: nil,
	}
	return g
}
func (g *game) Update() (err error) {
	args := g.Scene.Update()
	switch args := args.(type) {
	case error:
		err = args
	}
	return
}
func (g *game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(draws.Image{Image: screen})
}
func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scene.ScreenSizeX, scene.ScreenSizeY
}
