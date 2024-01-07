package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
)

// Separate main function from game.go to avoid collapsing the game package.
// Music directory is moved to cmd/gosu, because it is not a part of the game logic.
func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	root := os.DirFS(dir)

	g := gosu.NewGame(root)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
