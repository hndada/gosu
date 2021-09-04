package main

import (
	"log"

	"net/http"
	_ "net/http/pprof"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	g := gosu.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
