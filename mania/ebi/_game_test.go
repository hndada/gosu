package ebi

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"testing"
)

func TestGamePlay(t *testing.T) {
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("gosu!")
	// c := mania.NewChart("../test2.osu")
	// g := NewGame(c)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
