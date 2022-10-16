package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
)

// //go:embed skin
// var skin embed.FS

// //go:embed music
// var music embed.FS

func main() {
	g := gosu.NewGame([]gosu.ModeProp{piano.ModePiano4, piano.ModePiano7, drum.ModeDrum}) //, skin, music)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
