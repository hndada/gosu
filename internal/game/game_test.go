package game

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	g := NewGame()
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(g.ScreenWidth, g.ScreenHeight) // fixed in prototype
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetMaxTPS(g.MaxTPS)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}
