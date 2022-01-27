package mania

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/engine/ui"
)

const maxScore = 1e6

// Theorem: LNTail can't be unscored when key state is press or idle.
func (s *Scene) judge(e keyEvent) {
	i := s.staged[e.Key] // index of a staged note
	if i < 0 {
		return
	}
	n := s.chart.Notes[i] // staged note
	// if n.playSE != nil {
	// 	n.playSE()
	// }
	keyAction := KeyAction(s.lastPressed[e.Key], e.Pressed)
	timeDiff := n.Time - e.Time

	judgeable := func(t common.NoteType, keyAction int) bool {
		if t == TypeLNTail {
			return keyAction == release
		}
		return keyAction == press
	}
	judge := func(t common.NoteType, keyAction int, timeDiff int64) common.Judgment {
		if !judgeable(t, keyAction) {
			return empty
		}
		if timeDiff < 0 { // absolute value
			timeDiff *= -1
		}
		for _, j := range Judgments {
			if timeDiff <= j.Window {
				return j
			}
		}
		return empty // When a player hit too early
	}
	j := judge(n.Type, keyAction, timeDiff)
	s.applyScore(i, j)

	// if timeDiff <= Miss.Window {
	// 	ts := s.jm.NewTimingSprite(timeDiff)
	// 	s.timingSprites = append(s.timingSprites, ts)
	// }
}

func (s *Scene) applyScore(i int, j common.Judgment) {
	n := s.chart.Notes[i]
	if j == empty || n.scored {
		return
	}
	s.chart.Notes[i].scored = true
	s.chart.Notes[i].Sprite.Saturation = 0.5
	s.chart.Notes[i].Sprite.Dimness = 0.3
	s.staged[n.key] = n.next

	for idx, j2 := range Judgments {
		if j == j2 {
			s.judgeCounts[idx]++
			break
		}
	}
	switch common.Settings.ScoreMode {
	case common.ScoreModeNaive:
		unit := maxScore / float64(len(s.chart.Notes))
		s.score += unit * j.Value
		if s.hp > 0 {
			s.hp += 0.005 * j.Value
			if s.hp > 100 {
				s.hp = 100
			} else if s.hp < 0 {
				s.hp = 0
			}
		}
	case common.ScoreModeWeighted:
		// score
		if j.Value == 0 {
			s.score += math.Max(-800, -4*n.score) // not lower than -800
			if s.score < 0 {                      // score is non-negative
				s.score = 0
			}
		} else {
			s.score += n.score * j.Value * (1 + s.karma/100) * 0.5
		}

		// karma
		if j.Penalty == 0 {
			s.karma += n.karma
			if s.karma > 100 {
				s.karma = 100
			}
		} else {
			s.karma -= j.Penalty
			if s.karma < 0 {
				s.karma = 0
			}
		}

		// hp
		if s.hp > 0 {
			s.hp += n.hp * j.HP
			if s.hp > 100 {
				s.hp = 100
			} else if s.hp < 0 {
				s.hp = 0
			}
		}
	}

	// combo
	if j != Miss {
		s.combo++
	} else {
		s.combo = 0
	}

	if n.Type != TypeLNTail && j != Miss {
		s.playSE()
		// if n.playSE != nil {
		// 	n.playSE()
		// } else {
		// 	s.playSE() // default sample effect
		// }
	}
	for idx, j2 := range Judgments {
		if j == j2 {
			s.judgeSprite[idx].Rep = 2 // temp
			s.judgeSprite[idx].BornTime = time.Now()
			break
		}
	}

	switch n.Type {
	case TypeLNTail:
		s.LightingLN[n.key].Rep = 0
	}
	if j != Miss {
		switch n.Type {
		case typeNote:
			s.Lighting[n.key].BornTime = time.Now()
			s.Lighting[n.key].Rep = 1
		case TypeLNHead:
			s.LightingLN[n.key].Rep = ui.RepInfinite
		}
		// apply one more for LNTail when LNHead is missed
		if n.Type == TypeLNHead && j == Miss {
			s.applyScore(n.next, Miss)
		}
	}
}

func (s *Scene) drawCombo(screen *ebiten.Image) {
	gap := int(common.Settings.ComboGap * common.DisplayScale())
	str := fmt.Sprint(s.combo)
	w := s.combos[0].W
	wNumbers := (w-gap)*len(str) + gap
	for i, letter := range str {
		num, _ := strconv.ParseInt(string(letter), 0, 0)
		s.combos[num].X = common.Settings.ScreenSize.X/2 - wNumbers/2 + i*(w-gap)
		s.combos[num].Draw(screen)
	}
}

func (s *Scene) drawScore(screen *ebiten.Image) {
	str := fmt.Sprintf("%.0f", s.score)
	w := s.scores[0].W
	wNumbers := w * len(str)
	for i, letter := range str {
		num, _ := strconv.ParseInt(string(letter), 0, 0)
		s.scores[num].X = common.Settings.ScreenSize.X - wNumbers + i*w
		s.scores[num].Draw(screen)
	}
}

// func (s Score) JudgeCounts() []int { return s.Counts[:] }
// func (s Score) IsFullCombo() bool  { return s.Counts[4] == 0 }
// func (s Score) IsPerfect() bool {
// 	for _, c := range s.Counts[2:] {
// 		if c != 0 {
// 			return false
// 		}
// 	}
// 	return true
// }
