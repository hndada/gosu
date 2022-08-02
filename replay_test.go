package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hndada/gosu/parse/osr"
)

func TestReplayScore(t *testing.T) {
	c, err := NewChart("7k-gt.osu")
	if err != nil {
		panic(err)
	}
	rd, err := ioutil.ReadFile("7k-gt.osr")
	if err != nil {
		panic(err)
	}
	rf, err := osr.Parse(rd)
	if err != nil {
		panic(err)
	}
	s := NewScenePlay(c)
	s.ReplayMode = true
	s.ReplayStates = ExtractReplayState(rf, c.Parameter.KeyCount)
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
	fmt.Printf("New score: %d/1.1m\n", s.CurrentScore())
}
