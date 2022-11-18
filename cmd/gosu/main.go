package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
)

func main() {
	g := gosu.NewGame(os.DirFS(""))
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
