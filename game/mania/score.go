package mania

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
)

const maxScore = 1e6

// LNTail 이면서 unscored이고 press나 idle일순 없음
func (s *Scene) judge(e keyEvent) {
	i := s.staged[e.Key] // index of a staged note
	if i < 0 {
		return // todo: play sfx
	}
	n := s.chart.Notes[i] // staged note
	keyAction := KeyAction(s.lastPressed[e.Key], e.Pressed)
	timeDiff := n.Time - e.Time

	// Idle, Hit, Release, Hold (lost는 아예 별개의 개념. 시간 지나도록 X면)
	// 일반   노트: X, O, X, X
	// 롱노트 머리: X, O, X, X (단, miss시 꼬리까지 miss)
	// 롱노트 꼬리: X, X, O, X (현재 hold 시 HP보너스 생략)
	judgeable := func(t game.NoteType, keyAction int) bool { // judge 가능 action이나 premature(너무 빨리 누른 경우)인 경우 score 안됨.
		if t == TypeLNTail {
			return keyAction == release
		}
		return keyAction == press
	}
	judge := func(t game.NoteType, keyAction int, timeDiff int64) game.Judgment {
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
		return empty // 너무 빨리 누름. 너무 늦게 누른 경우(아예 안 누르다)는 scene update에서 별도 처리
	}
	j := judge(n.Type, keyAction, timeDiff)
	s.applyScore(i, j)

	// if timeDiff <= Miss.Window {
	// 	ts := s.jm.NewTimingSprite(timeDiff)
	// 	s.timingSprites = append(s.timingSprites, ts)
	// }
}

// LNTail은 롱노트 끝나기 전까지 계속 staged. 처음 scored 된 뒤로는 score 영향 안 끼침
func (s *Scene) applyScore(i int, j game.Judgment) {
	n := s.chart.Notes[i]
	if j == empty || n.scored { // scored되었는데 judge될 대상은 미리 뗀 LNTail 밖에 없음. LNTail은 scene update에서 별도 처리
		return
	}
	s.chart.Notes[i].scored = true
	s.chart.Notes[i].Sprite.Saturation = 0.5
	s.chart.Notes[i].Sprite.Dimness = 0.3
	s.staged[n.Key] = n.next

	for idx, j2 := range Judgments {
		if j == j2 {
			s.judgeCounts[idx]++
			break
		}
	}

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
	// combo
	if j != Miss {
		s.combo++
	} else {
		s.combo = 0
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
	if n.Type != TypeLNTail && j != Miss {
		s.playSE()
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
		s.LightingLN[n.Key].Rep = 0
	}
	if j != Miss {
		switch n.Type {
		case typeNote:
			s.Lighting[n.Key].BornTime = time.Now()
			s.Lighting[n.Key].Rep = 1
		case TypeLNHead:
			s.LightingLN[n.Key].Rep = game.RepInfinite
		}
		// apply one more for LNTail when LNHead is missed
		if n.Type == TypeLNHead && j == Miss {
			s.applyScore(n.next, Miss)
		}
	}
}

func (s *Scene) drawCombo(screen *ebiten.Image) {
	gap := int(game.Settings.ComboGap * game.DisplayScale())
	str := fmt.Sprint(s.combo)
	w := s.combos[0].W
	wNumbers := (w-gap)*len(str) + gap
	for i, letter := range str {
		num, _ := strconv.ParseInt(string(letter), 0, 0)
		s.combos[num].X = game.Settings.ScreenSize.X/2 - wNumbers/2 + i*(w-gap)
		s.combos[num].Draw(screen)
	}
}

func (s *Scene) drawScore(screen *ebiten.Image) {
	str := fmt.Sprintf("%.0f", s.score)
	w := s.scores[0].W
	wNumbers := w * len(str)
	for i, letter := range str {
		num, _ := strconv.ParseInt(string(letter), 0, 0)
		s.scores[num].X = game.Settings.ScreenSize.X - wNumbers + i*w
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
