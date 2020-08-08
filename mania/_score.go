package mania

import (
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/tools"
	"math"
)

// mania beatmap loader
// legacy replay 분석 좀 더, 이후 test 여부 결정
// (키보드 input: eiannone/keyboard module)

// 놓친 롱노트 끝날 때 리플레이가 어떻게 박히는지는 아직 확인 안함
// 시간 내에 correct한 action이 없을 경우 마지막에 miss 판정 내고 끝나는 걸로 상정 -> 여러 케이스 확인해봐야함 (sv2)
// 1. 노트, 롱노트 미리 누른 상태로 안 떼고 있을 경우
// 2. 롱노트 처음에 잘 치다가 뗀 경우, 그리고 다시 친 경우
// 그런데, 계속 누르고 있었으면 그냥 1 1 로 지속됐을 것 같음
// 현재 missed, 등호포함 부등호인데 legacy 할때는 풀어야 할 수도 있음 -> 너무 귀찮아지면 그냥 생략
const MaxScore = 1e6
const holdUnitHP = 0.002 // 롱노트를 눌렀을 때 1ms 당 차오르는 체력
const (
	idle = iota
	press
	release
	hold
)

func KeyAction(last, now bool) int { // action are realized with 2 snapshots
	switch {
	case !last && !now:
		return idle
	case !last && now:
		return press
	case last && !now:
		return release
	case last && now:
		return hold
	default:
		panic("not reach")
	}
}

func readPressed(x, keymode int) []bool {
	mask := 1
	pressed := make([]bool, keymode)
	for i := 0; i < keymode; i++ {
		pressed[i] = x&mask != 0
		mask = mask << 1
	}
	return pressed
}

type PlayNote struct {
	Note
	Score       float64
	KarmaScore  float64
	HPScore     float64
	NextNoteIdx int
	Scored      bool // 판정 완료 여부. 명암 변화 및 tail에선 점수에까지 영향
	// Position
}

func loadPlayNotes(ns []Note, keymode int) []PlayNote {
	// 사실 strain 아니고 aggregate인데 일단 이대로
	pns := make([]PlayNote, len(ns))
	var totalStrain float64
	var idxQueues = make([][]int, keymode)
	for k := range idxQueues {
		idxQueues[k] = make([]int, 0, len(ns)/(keymode-1))
	}

	for i, n := range ns {
		pns[i].Note = n
		totalStrain += n.Strain
		idxQueues[n.Key] = append(idxQueues[n.Key], i)
	}

	var idxQueueCursors = make([]int, keymode)
	avgStrain := totalStrain / float64(len(ns))
	for i, pn := range pns {
		pns[i].Score = MaxScore * (pn.Strain / totalStrain)
		pns[i].KarmaScore = 1 * math.Min(pn.Strain/avgStrain, 2.5)      // 0 ~ 2.5
		pns[i].HPScore = 1 * math.Min(pn.Strain/(3*avgStrain)+2/3, 1.5) // 0 ~ 1.5
		pns[i].NextNoteIdx = idxQueues[pn.Key][idxQueueCursors[pn.Key]]
		idxQueueCursors[pn.Key]++
	}
	return pns
}

// 실제 play, keyboard input을 받으면 그때마다 staged 체크
func ProcessScore(bpath, rpath string) {
	b := ManiaBeatmap(bpath)
	// if err != nil {
	// 	panic(err)
	// }
	r := game.ParseOsuReplay(rpath)
	pns := loadPlayNotes(b.Notes, b.Keymode)

	var keymode int // keymode := int(b.Difficulty["CircleSize"])

	var time int64
	var totalScore float64
	var karma float64 = 100
	var hp float64 = 100
	var staged = make([]int, keymode)
	var pressed = make([]bool, keymode)
	var lastPressed = make([]bool, keymode)

	inRange := func(elapsed int64) bool { return elapsed <= Judgements[4].Window }
	scorable := func(n PlayNote, action int) bool {
		switch n.NoteType {
		case NtHoldTail:
			return action == release
		default:
			return action == press
		}
	}
	missed := func(elapsed int64) bool { return elapsed > Judgements[3].Window }

	for _, ra := range r.ReplayData {
		time += ra.W
		pressed = readPressed(int(ra.X), keymode)

		for k, i := range staged {
			n := pns[i]
			ka := KeyAction(lastPressed[k], pressed[k])
			if inRange(n.Time-time) && scorable(n, ka) || missed(time-n.Time) { // including late miss
				if !pns[i].Scored {
					totalScore += score(time, n, &karma, &hp)
				}
				staged[k] = n.NextNoteIdx
			} else if n.NoteType == NtHoldTail {
				if !pressed[k] && !pns[i].Scored {
					totalScore += score(time, n, &karma, &hp) // when u release ln too early
					pns[i].Scored = true
				}
				if !lastPressed[k] {
					hp -= 4 * holdUnitHP * float64(ra.W)
				} else {
					hp += holdUnitHP * float64(ra.W)
				}
			}
		}
		if n := copy(lastPressed, pressed); n != keymode {
			panic("copy failed")
		}
	}
}

func score(time int64, n PlayNote, ptrKarma, ptrHP *float64) (s float64) {
	karma, hp := *ptrKarma, *ptrHP
	judge := func(diff int64) Judgement {
		absTime := tools.AbsInt(time)
		for _, judge := range Judgements {
			if absTime <= judge.Window {
				return judge
			}
		}
		return Judgements[4]
		// panic("not reach")
	}(n.Time - time)

	if judge.Value == 0 {
		s = math.Min(-1e4, -4*n.Score) // not lower than -10,000
	} else {
		s = n.Score * judge.Value * (1 + karma/100) * 0.5
	}
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
	hp += n.HPScore * judge.HP

	*ptrKarma = karma
	*ptrHP = hp
	return
}

// func fetchNextNote(start int, notes []Note) int {
// 	key := notes[start].Key
// 	for i, n := range notes[start+1:] {
// 		if n.Key == key {
// 			return start + 1 + i
// 		}
// 	}
// 	return -1
// }
