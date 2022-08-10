package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/parse/osr"
	"github.com/hndada/gosu/parse/osu"
)

// This is for testing parsing replay and simulating playing.
func main() {
	cpath := "music/doppelganger/LeaF - Doppelganger (Jinjin) [jakads' Extra].osu"
	b, err := os.ReadFile(cpath)
	if err != nil {
		panic(err)
	}
	o, err := osu.Parse(b)
	if err != nil {
		panic(err)
	}
	c, err := gosu.NewChartFromOsu(o)
	if err != nil {
		panic(err)
	}
	// rd, err := os.ReadFile("replay/osu!topus! - nekodex - circles! [Hard] (2022-08-10) OsuMania.osr")
	rd, err := os.ReadFile("replay/replay-mania_1023967_492000477.osr")
	if err != nil {
		panic(err)
	}
	rf, err := osr.Parse(rd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Song: %s\n", c.MusicName)
	fmt.Printf("Diff: %s\n", c.ChartName)
	fmt.Printf("Player: %s\n", rf.PlayerName)
	fmt.Printf("Original score: %d/1m\n", rf.Score)
	fmt.Printf("New score: %d/1.1m\n", gosu.ReplayScore(c, rf))
	g := gosu.NewGame()
	g.Scene = gosu.NewScenePlay(c, cpath, rf, true)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
