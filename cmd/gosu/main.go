package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
)

func main() {
	g := gosu.NewGame([]gosu.ModeProp{piano.ModePiano4, piano.ModePiano7, drum.ModeDrum})
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
