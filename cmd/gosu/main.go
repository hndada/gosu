package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu"
	"log"
)

func main() {
	g := gosu.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
