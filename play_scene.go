package main

import "github.com/hajimehoshi/ebiten/v2"

// ScenePlay: struct, PlayScene: function
type ScenePlay struct {
	// Music Streamer
	Tick int

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
}

func NewScenePlay(c *Chart) *ScenePlay {
	s := new(ScenePlay)
	s.Tick = -2 * ebiten.MaxTPS() // Put 2 seconds of waiting
	s.PlayNotes = NewPlayNotes(c) // Todo: add Mods to input param
	s.KeySettings = KeySettings[s.Chart.Parameter.KeyCount]
	return s
}

func (s *ScenePlay) Update() {
	s.Tick++
	for k, p := range s.Pressed {
		s.LastPressed[k] = p
		s.Pressed[k] = ebiten.IsKeyPressed(s.KeySettings[k])
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
		if j := Verdict(n, td, s.KeyAction(n.Key)); j.Window != 0 {
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
	return int64(float64(s.Tick) / float64(ebiten.MaxTPS()))
}

func (s ScenePlay) IsFinished() bool {
	return s.Time() > 3000+s.PlayNotes[len(s.PlayNotes)-1].Time
}
