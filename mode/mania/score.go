package mania

import (
	"github.com/hndada/gosu/mode"
	"math"
)

// 핵심은, 롱노트를 놔서 최종 미스 판정을 받았더라도 staged에 ln tail 이 있어야 한다는 것
// Tail 이면서 unscored이고 press나 idle일순 없음
// 안쳐서 미스난 건 event가 없을 수 있음, 바깥에서 진행되어야 함
// todo: 리플레이 -> 키보드 인풋처럼
// 리플레이 구조: 마지막 status 시간, 레이아웃 키state

type Score struct {
	mode.BaseScore
	Counts [len(judgments)]int
}

func (s Score) JudgeCounts() []int { return s.Counts[:] }
func (s Score) IsFullCombo() bool  { return s.Counts[4] == 0 }
func (s Score) IsPerfect() bool {
	for _, c := range s.Counts[2:] {
		if c != 0 {
			return false
		}
	}
	return true
}

const holdUnitHP = 0.002 // 롱노트를 눌렀을 때 1ms 당 차오르는 체력

// todo: CalcLV 에 추가
func (c *Chart) CalcScore2() {
	var sumStrain float64
	for _, n := range c.Notes {
		sumStrain += n.strain
	}
	var avgStrain float64
	if len(c.Notes) != 0 {
		avgStrain = sumStrain / float64(len(c.Notes))
	}
	for i, n := range c.Notes {
		c.Notes[i].score = mode.MaxScore * (n.strain / sumStrain)
		c.Notes[i].karma = math.Min(n.strain/avgStrain, 2.5)      // 0 ~ 2.5
		c.Notes[i].hp = math.Min(n.strain/(3*avgStrain)+2/3, 1.5) // 0 ~ 1.5
	}
}

// 놓친 롱노트 끝날 때 리플레이가 어떻게 박히는지는 아직 확인 안함
// 시간 내에 correct한 action이 없을 경우 마지막에 miss 판정 내고 끝나는 걸로 상정 -> 여러 케이스 확인해봐야함 (sv2)
// 1. 노트, 롱노트 미리 누른 상태로 안 떼고 있을 경우
// 2. 롱노트 처음에 잘 치다가 뗀 경우, 그리고 다시 친 경우
// 그런데, 계속 누르고 있었으면 그냥 1 1 로 지속됐을 것 같음
// 현재 missed, 등호포함 부등호인데 legacy 할때는 풀어야 할 수도 있음 -> 너무 귀찮아지면 그냥 생략

// 실제 play, keyboard input을 받으면 그때마다 staged 체크
func inRange(time int64) bool           { return time <= miss.Window }
func lost(time int64) bool              { return time < -bad.Window }
func drainable(n Note, time int64) bool { return n.scored && time < miss.Window }
func scoreable(n Note, action int, time int64) bool {
	if n.Type == TypeLNTail {
		return action == release
	}
	return inRange(time) && action == press
}
// func outer() { // 이미 n.scored == false, 검사할 필요 없음
// 	if lost(t) {
// 		s.updateScore(i, miss)
// 		if n.Type == TypeLNHead {
// 			s.updateScore(n.next, miss) // todo: scored 에 true대입하는걸 updateScore에서 하기
// 		}
// 	}
// }

func judge(n Note, action int, time int64) (mode.Judgment, bool) { // bool: judged
	if !scoreable(n, action, time) {
		return mode.Judgment{}, false
	}
	if time < 0 {
		time *= -1
	}
	for _, j := range judgments {
		if time <= j.Window {
			return j, true
		}
	}
	return miss, true // reaches only when release ln too early
}
func (s *Scene) processScore(e mode.KeyEvent) {
	i := s.staged[e.KeyCode]
	n := s.chart.Notes[i]
	lastPressed := lastLog[e.Key]
	action := KeyAction(lastPressed, e.State)
	time := n.Time - e.Time

	if j, judged := judge(n, action, time); judged && !n.scored {
		s.updateScore(i, j)
	}
	if drainable(s.chart.Notes[i], time) { // scored가 이미 돼있을 수도 있어서 분리
		s.staged[e.Key] = n.next
	}

	if n.Type == TypeLNTail {
		var holdTime float64
		if e.Time > n.Time {
			holdTime = float64(n.Time - lastLog.Time)
		} else {
			holdTime = float64(e.Time - lastLog.Time)
		}
		if holdTime < 0 {
			holdTime = 0
		}
		switch action {
		case hold, release: // release: 1, 0
			s.hp += holdUnitHP * holdTime
		case idle, press: // press: 0, 1
			s.hp -= 4 * holdUnitHP * holdTime
		}
	}
}

func (s *Scene) updateScore(i int, j mode.Judgment) {
	// last pressed: s.logs[len(logTime)-1]
	// if time is same (누락된 것)-> 그 이전걸로.
	n := s.chart.Notes[i]
	if j.Value == 0 {
		s.score += math.Min(-1e4, -4*n.score) // not lower than -10,000
	} else {
		s.score += n.score * j.Value * (1 + s.karma/100) * 0.5
	}
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
	s.hp += n.hp * j.HP
	s.chart.Notes[i].scored = true
}
