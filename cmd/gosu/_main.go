package main

// This is for testing parsing replay and simulating playing.
import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu"
	"github.com/hndada/gosu/parse/osr"
	"github.com/hndada/gosu/parse/osu"
)

func main() {
	cpath := "music/circles/nekodex - circles! (MuangMuangE) [Hard].osu"
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
	rd, err := os.ReadFile("replay/circles.osr")
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
	// fmt.Printf("New score: %.0f/1.1m\n", s.CurrentScore())
	g := gosu.NewGame()
	g.Scene = gosu.NewScenePlay(c, cpath, rf)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
