package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
)

func main() {
	g := mode.NewGame([]mode.Mode{piano.ModePiano4, piano.ModePiano7})
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
