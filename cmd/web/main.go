package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mode/drum"
	"github.com/hndada/gosu/game/mode/piano"
)

// //go:embed skin
// var skin embed.FS

// //go:embed music
// var music embed.FS

func main() {
	g := game.NewGame([]game.ModeProp{piano.ModePiano4, piano.ModePiano7, drum.ModeDrum}) //, skin, music)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
