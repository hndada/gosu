package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
)

func _main() {
	g := gosu.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
