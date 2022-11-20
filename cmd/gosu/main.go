package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	g := gosu.NewGame(os.DirFS(dir))
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
