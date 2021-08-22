package mania

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
)

const maxScore = 1e6

// type Score struct {
// 	game.BaseScore
// 	Counts [len(Judgments)]int
// }

// 'Late hit' gets negative value in hit time difference
// 놓친 롱노트 끝날 때 리플레이가 어떻게 박히는지는 아직 확인 안함
// 시간 내에 correct한 action이 없을 경우 마지막에 miss 판정 내고 끝나는 걸로 상정 -> 여러 케이스 확인해봐야함 (sv2)
// 1. 노트, 롱노트 미리 누른 상태로 안 떼고 있을 경우
// 2. 롱노트 처음에 잘 치다가 뗀 경우, 그리고 다시 친 경우
// 그런데, 계속 누르고 있었으면 그냥 1 1 로 지속됐을 것 같음
// 현재 missed, 등호포함 부등호인데 legacy 할때는 풀어야 할 수도 있음 -> 너무 귀찮아지면 그냥 생략

// 핵심은, 롱노트를 놔서 최종 미스 판정을 받았더라도 staged에 ln tail 이 있어야 한다는 것
// Tail 이면서 unscored이고 press나 idle일순 없음

// LN holdHP. 추가하고자 한다면 lastPressed에 time value까지 추가
// log 삭제. 리플레이는 별도 함수에서 처리하는 걸로
func (c *Chart) allotScore() {
	var sumStrain float64
	for _, n := range c.Notes {
		sumStrain += n.strain
	}
	var avgStrain float64
	if len(c.Notes) != 0 {
		avgStrain = sumStrain / float64(len(c.Notes))
	}
	for i := range c.Notes {
		n := c.Notes[i]
		c.Notes[i].score = maxScore * (n.strain / sumStrain)
		c.Notes[i].karma = math.Min(n.strain/avgStrain, 2.5)          // 0 ~ 2.5
		c.Notes[i].hp = math.Min(n.strain/(3*avgStrain)+2.0/3.0, 1.5) // 0 ~ 1.5
	}
}

func (s *Scene) judge(e keyEvent) {
	i := s.staged[e.key] // index of a staged note
	if i < 0 {
		return // todo: play sfx
	}
	n := s.chart.Notes[i] // staged note
	keyAction := KeyAction(s.lastPressed[e.key].Value, e.pressed)
	timeDiff := n.Time - e.time

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
	s.timeDiffs = append(s.timeDiffs, timeDiff)
	s.applyScore(i, j)
	s.lastPressed[e.key] = TimeBool{Time: e.time, Value: e.pressed}
}

// LNTail은 롱노트 끝나기 전까지 계속 staged. 처음 scored 된 뒤로는 score 영향 안 끼침
func (s *Scene) applyScore(i int, j game.Judgment) {
	n := s.chart.Notes[i]
	if j == empty || n.scored { // scored되었는데 judge될 대상은 미리 뗀 LNTail 밖에 없음. LNTail은 scene update에서 별도 처리
		return
	}
	s.chart.Notes[i].scored = true
	s.staged[n.Key] = n.next

	switch j {
	case Kool:
		s.judgeCounts[0]++
	case Cool:
		s.judgeCounts[1]++
	case Good:
		s.judgeCounts[2]++
	case Bad:
		s.judgeCounts[3]++
	case Miss:
		s.judgeCounts[4]++
	default:
		panic("not reach")
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
	if s.lastJudge.Penalty < j.Penalty {
		s.lastJudge = j
	}

	// apply one more for LNTail when LNHead is missed
	if n.Type == TypeLNHead && j == Miss {
		s.applyScore(n.next, Miss)
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

// const holdUnitHP = 0.002 // 롱노트를 눌렀을 때 1ms 당 차오르는 체력

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
