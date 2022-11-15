package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/game"
)

func main() {
	g := game.NewGame(os.DirFS(""))
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
