package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu"
)

func main() {
	g := gosu.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
