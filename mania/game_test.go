package mania

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/mania/ebi"
	"log"
	"testing"
)

func TestGamePlay(t *testing.T) {
	ebiten.SetWindowSize(mania.ScreenWidth, mania.ScreenHeight)
	ebiten.SetWindowTitle("gosu!")
	c := NewChart("test2.osu")
	g := mania.NewGame(c)
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}
