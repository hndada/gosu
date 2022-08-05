package main

import (
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

	Speed float64
	TransPoint
}

// var MaxTPS int64 = 1000 // Todo: mapping to ebiten

func SecondToTick(t int64) int64 { return t * MaxTPS }
func NewScenePlay(c *Chart) *ScenePlay {
	s := new(ScenePlay)
	s.Tick = -2 * MaxTPS // Put 2 seconds of waiting
	s.Chart = c
	s.PlayNotes, s.StagedNotes = NewPlayNotes(c) // Todo: add Mods to input param
	// s.KeySettings = KeySettings[s.Chart.Parameter.KeyCount]
	s.JudgmentCounts = make([]int, 5)
	s.LastPressed = make([]bool, c.KeyCount)
	s.Pressed = make([]bool, c.KeyCount)
	s.Karma = 1
	s.TransPoint = TransPoint{
		s.Chart.SpeedFactors[0],
		s.Chart.Tempos[0],
		s.Chart.Volumes[0],
		s.Chart.Effects[0],
	}
	return s
}

func (s *ScenePlay) Update() {
	s.Tick++
	for s.Time() < s.SpeedFactor.Next.Time {
		s.SpeedFactor = s.SpeedFactor.Next
	}
	for s.Time() < s.Tempo.Next.Time {
		s.Tempo = s.Tempo.Next
	}
	for s.Time() < s.Volume.Next.Time {
		s.Volume = s.Volume.Next
	}
	for s.Time() < s.Effect.Next.Time {
		s.Effect = s.Effect.Next
	}

	for k, p := range s.Pressed {
		s.LastPressed[k] = p
		if s.ReplayMode {
			for s.ReplayCursor < len(s.ReplayStates)-1 && s.Time() > s.ReplayStates[s.ReplayCursor].Time {
				s.ReplayCursor++
			}
			s.ReplayCursor--
			if s.ReplayCursor < 0 {
				s.ReplayCursor = 0
			}
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
}

func (s *ScenePlay) Draw() {
	if s.IsFinished() {
		s.DrawClear()
		return
	}
	s.DrawBG()
	s.DrawField()
	s.DrawLongNote()
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
