package mania

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/beatmap"
	"github.com/hndada/gosu/tools"
)

const MaxScore = 1e6

func readPressed(x, keymode int) []bool {
	mask := 1
	pressed := make([]bool, keymode)
	for i := 0; i < keymode; i++ {
		pressed[i] = x&mask != 0
		mask = mask << 1
	}
	return pressed
}

const (
	idle = iota
	hit
	release
	hold
)

// 언제나 구간으로 실현
func KeyAction(last, now bool) int {
	switch {
	case !last && !now:
		return idle
	case !last && now:
		return hit
	case last && !now:
		return release
	case last && now:
		return hold
	default:
		panic("not reach")
	}
}

// 실제 play,
// keyboard input을 받으면 그때마다 staged 체크
// 범위 내면 score 함수 호출, 체력 함수 호출
func ProcessScore(bpath, rpath string) {
	b, err := beatmap.ParseBeatmap(bpath)
	if err != nil {
		panic(err)
	}
	r := game.ReadLegacyReplay(rpath)
	var keymode int // keymode := int(b.Difficulty["CircleSize"])

	var time int64
	var staged = make([]int, keymode)
	var pressed = make([]bool, keymode)
	var lastPressed = make([]bool, keymode)
	var unitScore, totalScore float64
	var karma float64 = 100

	// 놓친 롱노트 끝날 때 리플레이가 어떻게 박히는지는 아직 확인 안함
	// 시간 내에 correct한 action이 없을 경우 마지막에 miss 판정 내고 끝나는 걸로 상정 -> 여러 케이스 확인해봐야함 (sv2)
	// 1. 노트, 롱노트 미리 누른 상태로 안 떼고 있을 경우
	// 2. 롱노트 처음에 잘 치다가 뗀 경우, 그리고 다시 친 경우
	// 그런데, 계속 누르고 있었으면 그냥 1 1 로 지속됐을 것 같음

	// && time-n.Time < Judgements[3].Window // score에서 알아서 걸러짐
	inRange := func(elapsed int64) bool { return elapsed <= Judgements[4].Window }
	scorable := func(n PlayingNote, action int) bool {
		switch n.NoteType {
		case NtHoldTail:
			return action == release
		default:
			return action == hit
		}
	}
	missed := func(elapsed int64) bool { return elapsed > Judgements[3].Window }

	for _, ra := range r.ReplayData {
		time += ra.W
		pressed = readPressed(int(ra.X), keymode)

		for k, n := range staged {
			ka := KeyAction(lastPressed[k], pressed[k])
			if inRange(n.Time-time) && scorable(n, ka) || missed(time-n.Time) { // 자동 미스 포함
				unitScore, karma = score(time, n, karma)
				totalScore += unitScore
				hp()
				staged[k] = n.NextNote
			} else if n.NoteType == NtHoldTail {
				// 마크만 하고 hp 처리
				// 놓는 action일 경우 scored=true
			}
		}
		if n := copy(lastPressed, pressed); n != keymode {
			panic("copy failed")
		}
	}
}

// done 조건: 꼬리 노트 판정 받고 현재 time이 롱노트 꼬리 노트보다 크거나 같을때
type PlayingNote struct {
	Note
	Score      float64
	KarmaScore float64
	NextNote   *Note
	Scored     bool // 판정 완료 여부. 명암 변화 및 tail에선 점수에까지 영향
	// Position

	// idx        int
	// result Judgement
}

func score(time int64, n PlayingNote, karma float64) (float64, float64) {
	judge := func(diff int64) Judgement { // 시간 다 지나서 생긴 미스 판정도 처리됨
		absTime := tools.AbsInt(time)
		for _, judge := range Judgements {
			if absTime <= judge.Window {
				return judge
			}
		}
		panic("not reach")
	}(n.Time - time)

	s := n.Score * judge.Value * (1 + karma/100) * 0.5
	if judge.Penalty == 0 {
		karma += n.KarmaScore
		if karma > 100 {
			karma = 100
		}
	} else {
		karma -= judge.Penalty
		if karma < 0 {
			karma = 0
		}
	}
	return s, karma
}

// 처음 PlayingNote 만들 때 싹 로딩
func baseScore(n PlayingNote) float64 {
	return n.Strain / totalStrain * MaxScore // 사실 strain 아니고 aggregate인데 일단 이대로
}
func noteKarma(n PlayingNote) float64 { // struct 로 삽입
	return n.Strain / (totalStrain / len(notes)) * 1
}

// 처음 PlayingNote 만들 때 한꺼번에 로드
func fetchNextNote(start int, notes []Note) int {
	key := notes[start].Key
	for i, n := range notes[start+1:] {
		if n.Key == key {
			return start + 1 + i
		}
	}
	return -1
}
