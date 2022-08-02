package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

// ScenePlay: struct, PlayScene: function
type ScenePlay struct {
	// Music Streamer
	Tick int64

	Chart       *Chart
	PlayNotes   []*PlayNote
	KeySettings []ebiten.Key // Todo: separate ebiten

	Pressed     []bool
	LastPressed []bool

	StagedNotes    []*PlayNote
	Combo          int
	Karma          float64
	KarmaSum       float64
	JudgmentCounts []int

	// In dev
	ReplayMode   bool
	ReplayStates []ReplayState
	ReplayCursor int
}

var MaxTPS int64 = 1000 // Todo: mapping to ebiten

func SecondToTick(t int64) int64 { return t * MaxTPS }
func NewScenePlay(c *Chart) *ScenePlay {
	s := new(ScenePlay)
	s.Tick = -2 * MaxTPS // Put 2 seconds of waiting
	s.Chart = c
	s.PlayNotes = NewPlayNotes(c) // Todo: add Mods to input param
	// s.KeySettings = KeySettings[s.Chart.Parameter.KeyCount]
	s.JudgmentCounts = make([]int, 5)
	s.LastPressed = make([]bool, c.Parameter.KeyCount)
	s.Pressed = make([]bool, c.Parameter.KeyCount)
	return s
}

func (s *ScenePlay) Update() {
	s.Tick++
	for k, p := range s.Pressed {
		s.LastPressed[k] = p
		if s.ReplayMode {
			for s.ReplayCursor < len(s.ReplayStates)-1 && s.Time() < s.ReplayStates[s.ReplayCursor].Time {
				s.ReplayCursor++
			}
			fmt.Println(s.Time(), s.ReplayStates[s.ReplayCursor].Time)
			s.Pressed = s.ReplayStates[s.ReplayCursor].Pressed
		} else {
			s.Pressed[k] = ebiten.IsKeyPressed(s.KeySettings[k])
		}
	}
	for k, n := range s.StagedNotes {
		if n == nil {
			continue
		}
		if n.Type != Tail && s.KeyAction(k) == Hit {
			n.PlaySE()
		}

		td := n.Time - s.Time() // Time difference; negative values means late hit
		if n.Scored {
			if n.Type != Tail {
				panic("non-tail note has not flushed")
			}
			if td < Miss.Window { // Keep Tail being staged until nearly ends
				s.StagedNotes[n.Key] = n.Next
			}
			continue
		}
		if j := Verdict(n.Type, s.KeyAction(n.Key), td); j.Window != 0 {
			s.Score(n, j)
		}
	}
	for _, n := range s.PlayNotes {
		n.UpdateSprite()
	}
}

func (s *ScenePlay) Draw() {
	if s.IsFinished() {
		s.DrawClear()
		return
	}
	s.DrawBG()
	s.DrawField()
	s.DrawNote()
	s.DrawCombo()
	s.DrawJudgment()
	s.DrawOthers() // Score, judgment counts and other states
}

func (s ScenePlay) Time() int64 {
	return int64(float64(s.Tick) / float64(MaxTPS) * 1000)
}

func (s ScenePlay) IsFinished() bool {
	return s.Time() > 3000+s.PlayNotes[len(s.PlayNotes)-1].Time
}
