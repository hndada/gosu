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

var (
	Kool = common.Judgment{Value: 1.0000, Penalty: 0, HP: 0.75, Window: 16, ComboBreak: false}
	Cool = common.Judgment{Value: 0.9375, Penalty: 0, HP: 0.5, Window: 40, ComboBreak: false} // 15/16
	Good = common.Judgment{Value: 0.625, Penalty: 4, HP: 0.25, Window: 70, ComboBreak: false} // 10/16
	Bad  = common.Judgment{Value: 0.25, Penalty: 10, HP: 0, Window: 100, ComboBreak: false}
	Miss = common.Judgment{Value: 0, Penalty: 25, HP: -3, Window: 150, ComboBreak: true}
)

const maxScore = 1e6
const scoreMode = scoreModeWeighted
const (
	scoreModeNaive = iota
	scoreModeWeighted
	scoreModeOsuLegacy
)

type scores struct {
	common.Scores
	Karma float64 // pseudo-public field: for readability
}

func newScores() scores {
	var s scores
	s.HP = 100
	s.Karma = 100
	s.Judgments = []common.Judgment{Kool, Cool, Good, Bad, Miss}
	s.JudgeCounts = make([]int, 5)
	return s
}
func (s scores) judgable(t common.NoteType, a common.KeyActionState) bool {
	if t == TypeLNTail {
		return a == common.Release
	}
	return a == common.Press
}

// Key sound generates considerable lagging
// Theorem: LNTail can't be unscored when key state is press or idle.
func (s *Scene) applyScore(i int, j common.Judgment) {
	n := s.chart.Notes[i]
	if j.Window == 0 || n.scored {
		return
	}
	s.chart.Notes[i].scored = true
	s.chart.Notes[i].Sprite.Saturation = 0.5
	s.chart.Notes[i].Sprite.Dimness = 0.3
	s.staged[n.key] = n.next
	s.CountJudge(j)

	switch scoreMode {
	case scoreModeNaive:
		unit := maxScore / float64(len(s.chart.Notes))
		s.Score += unit * j.Value
		if s.HP > 0 {
			s.HP += 0.005 * j.Value
			if s.HP > 100 {
				s.HP = 100
			} else if s.HP < 0 {
				s.HP = 0
			}
		}
	case scoreModeWeighted:
		// score
		if j.Value == 0 {
			s.Score += math.Max(-800, -4*n.score) // not lower than -800
			if s.Score < 0 {                      // score is non-negative
				s.Score = 0
			}
		} else {
			s.Score += n.score * j.Value * (1 + s.Karma/100) * 0.5
		}

		// karma
		if j.Penalty == 0 {
			s.Karma += n.karma
			if s.Karma > 100 {
				s.Karma = 100
			}
		} else {
			s.Karma -= j.Penalty
			if s.Karma < 0 {
				s.Karma = 0
			}
		}

		// hp
		if s.HP > 0 {
			s.HP += n.hp * j.HP
			if s.HP > 100 {
				s.HP = 100
			} else if s.HP < 0 {
				s.HP = 0
			}
		}
	}
	s.CountCombo(j)
	if n.Type != TypeLNTail && j != Miss {
		// s.playSE() // TODO: Laggy
		// if n.playSE != nil {
		// 	n.playSE()
		// } else {
		// 	s.playSE() // default sample effect
		// }
	}
	for idx, j2 := range s.Judgments {
		if j == j2 {
			s.judgeSprite[idx].Rep = 2 // TEMP
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
	str := fmt.Sprint(s.Combo)
	w := s.combos[0].W
	wNumbers := (w-gap)*len(str) + gap
	for i, letter := range str {
		num, _ := strconv.ParseInt(string(letter), 0, 0)
		s.combos[num].X = common.Settings.ScreenSizeX/2 - wNumbers/2 + i*(w-gap)
		s.combos[num].Draw(screen)
	}
}

func (s *Scene) drawScore(screen *ebiten.Image) {
	str := fmt.Sprintf("%.0f", s.Score)
	w := s.sceneUI.scores[0].W
	wNumbers := w * len(str)
	for i, letter := range str {
		num, _ := strconv.ParseInt(string(letter), 0, 0)
		s.sceneUI.scores[num].X = common.Settings.ScreenSizeX - wNumbers + i*w
		s.sceneUI.scores[num].Draw(screen)
	}
}
func lost(t int64) bool                { return t < -Bad.Window }             // A note goes passed without being hit
func discardable(n Note, t int64) bool { return n.scored && t < Miss.Window } // Whether a staged note is able to be discard
