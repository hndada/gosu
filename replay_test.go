package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hndada/gosu/parse/osr"
	"github.com/hndada/gosu/parse/osu"
)

func TestReplayScore(t *testing.T) {
	b, err := os.ReadFile("7k-gt.osu")
	if err != nil {
		panic(err)
	}
	o, err := osu.Parse(b)
	if err != nil {
		panic(err)
	}
	c, err := NewChartFromOsu(o)
	if err != nil {
		panic(err)
	}
	rd, err := os.ReadFile("7k-gt.osr")
	if err != nil {
		panic(err)
	}
	rf, err := osr.Parse(rd)
	if err != nil {
		panic(err)
	}
	s := NewScenePlay(c, "7k-gt.osu") // Temporary chart path
	s.ReplayMode = true
	s.ReplayStates = ExtractReplayState(rf, c.KeyCount)
	s.Tick = -2 * MaxTPS
	for !s.IsFinished() {
		s.Update()
		// if s.Tick%1000 == 0 {
		// 	fmt.Println(s.Tick, s.CurrentScore(), s.JudgmentCounts, s.ReplayCursor, s.Pressed)
		// }
	}
	fmt.Printf("Song: %s\n", c.MusicName)
	fmt.Printf("Diff: %s\n", c.ChartName)
	fmt.Printf("Player: %s\n", rf.PlayerName)
	fmt.Printf("Original score: %d/1m\n", rf.Score)
	fmt.Printf("New score: %.0f/1.1m\n", s.CurrentScore())
}
