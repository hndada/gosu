package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
)

func main() {
	// g := gosu.NewGame()
	// if err := ebiten.RunGame(g); err != nil {
	// 	panic(err)
	// }
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	gosu.Load(os.DirFS(dir))
	g := gosu.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
