package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	g := gosu.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
