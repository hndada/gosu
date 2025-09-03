package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	root := os.DirFS(dir)

	g, err := NewGame(root)
	if err != nil {
		panic(err)
	}

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
