package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
)

func main() {
	g := gosu.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
