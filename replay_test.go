package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/hndada/gosu/parse/osr"
)

func TestReplayScore(t *testing.T) {
	c, err := NewChart("4k.osu")
	if err != nil {
		panic(err)
	}
	rd, err := ioutil.ReadFile("4k.osr")
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
		if s.Tick%300 == 0 {
			fmt.Println(s.Tick, s.CurrentScore(), s.JudgmentCounts, s.ReplayCursor, s.Pressed)
		}
	}
}
