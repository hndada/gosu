package game_test

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	g := gosu.NewGame()
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(g.ScreenWidth, g.ScreenHeight) // fixed in prototype
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetMaxTPS(g.MaxTPS)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}