package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/game"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	root := os.DirFS(dir)

	g, err := game.NewGame(root)
	if err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
