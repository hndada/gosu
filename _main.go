package main

import (
	"os"

	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mode/drum"
	"github.com/hndada/gosu/game/mode/piano"
)

func main() {
	// l, err := net.Listen("tcp", ":54125")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Listening on %s\n", l.Addr())
	// go http.Serve(l, nil)

	fsys := os.DirFS("skin")
	game.LoadGeneralSkin(fsys)
	piano.LoadSkin(fsys)
	drum.LoadSkin(fsys)
	g := game.NewGame([]game.NewScenePlay{piano.NewScenePlay, drum.NewScenePlay})
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
