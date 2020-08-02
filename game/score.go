package game

import "math"

// No Mods except Rate one, thus no need to separate HitBonus (someday)

const MaxScore = 1e6

type Judgement struct {
	Value      int
	BonusValue int
	Bonus      int
	Punishment int
}

var Judgements [6]Judgement

// only NoMod for simple test plz
func init() {
	Values := [6]int{320, 300, 200, 100, 50, 0}
	BonusValues := [6]int{32, 32, 16, 8, 4, 0}
	Bonuses := [6]int{2, 1, 0, 0, 0, 0}
	Punishments := [6]int{0, 0, 8, 24, 44, 200}
	for i := range Judgements {
		Judgements[i] = Judgement{
			Values[i], BonusValues[i], Bonuses[i], Punishments[i],
		}
	}
}

// 동일 판정이어도, 노트별 난이도가 다르기 때문에 gosu에서는 점수가 달라지는게 당연지사

// 노트 판정 조건: 마지막 state가 idle(off)이면서 판정 범위 내에 pressed(on).
// 롱노트 판정 조건: missed가 아니고 마지막 state가 on이면서 꼬리 판정 범위 내에서 unpressed(off)
// 롱노트 시작 이후, 꼬리 판정 범위 이전에 off가 발생하면 꼬리 miss 판정.

// 노트 판정은 즉시 처리 후 score queue에. 한번 판정 처리된 노트는 다시 보지 않는다.
// 아예 안 눌러 미스가 발생한 경우 미스 한계 범위 넘어가면서 판정 처리. 이외에는 대기.
// 안친지 97ms (OD0기준 '50' 경계선) 넘어간 시점에서 replay 기록 남김. (OD10에선 30ms 빨리 기록됨)

// staged 노트의 array가 있어야 할 것 같다
// 곧 칠 예정이거나 이미 쳤어야 하는 범위 이내 노트들 및 롱노트의 경우 안 끝난 노트

// 한번 루프 돌때마다 하는 거
// 리플레이 액션을 읽는다
// 현재 state로 리플레이 액션을 업데이트한다
// staged 스캔, 현재 및 이전 state들로 판정 처리->score queue에 업데이트.
// staged에서 done 된 노트를 업데이트
// 마지막으로 현재 state를 이전 state에 복사한다

// closure, function 단위 init정도로 보면 될듯 (초깃값 설정)

type PlayingNote struct {
	// done 조건: 꼬리 노트 판정 받고 현재 time이 롱노트 꼬리 노트보다 크거나 같을때
	Note
	idx    int
	played bool // 판정 완료 여부
	holdon bool
	// position
	result Judgement

}

func ProcessScore() {
	// read beatmap
	// read replay
	const keymode = beatmap.Keymode

	var time int
	var stagedNotes = make([]int, keymode)

	var pressed = make([]bool, keymode)
	var lastPressed = make([]bool, keymode)

	// 언제나 구간으로 실현
	idle := func(k int) bool { return !lastPressed[k] && !pressed[k] }
	hit := func(k int) bool { return !lastPressed[k] && pressed[k] }
	release := func(k int) bool { return lastPressed[k] && !pressed[k] }
	hold := func(k int) bool { return lastPressed[k] && pressed[k] }

	for _, action := replay.ReplayData {
		time += action.w
		pressed = readPressed(action.x, keymode)
		for k, n := range stagedNotes {
			// 안 쳐서 미스 난 것도 처리해야함
			switch n.NoteType {
			case Note, HoldHead:
				if n.Time-time < MaxRange && hit(k) {
					queueScore(n)
					stagedNotes[k] = fetchNextNote(n.idx, notes)
				} else if time-n.Time > Bad.Range && !hit(k) {
					queueScore(n)
					stagedNotes[k] = fetchNextNote(n.idx, notes)
				}
				// 홀드 tail 까지 판정 완료되었어도 눌러도 회복 x (사실 이미 notes에 없다)
				// 홀드 tail도 처음에는 체력 감소가 있게 하자
			case HoldTail:
				switch {
				case hold(k): // 롱놋 켜져있을 때에는 언제나 회복. 꺼지고 미스 한계까지 누르고 있으면 미스.
					if time < n.Time {
						heal()
					} else if time > n.Time+Bad.Range {
						score()
					}
				case release(k): // 꼬리 판정 전에 놓으면 판정. 단, next를 fetch 하진 않음.
					if time < n.Time {
						heal()
					} else if time-n.Time < -Bad.Range {
						score() // 이때는 fetch 함
					}
				}

				// 시작 중간 끝
				if n.Time-time < MaxRange && hit(k) {
					queueScore(n)
					stagedNotes[k] = fetchNextNote(n.idx, notes) // 끝노트 fetch 할수도 있음
				} else if time-n.Time > Bad.Range {
					if release(k) {
						stagedNotes[k].played = true // queueScore에서 처리
						// 지난 시간 만큼 감소
					} else {
						// 지난 시간 만큼
					}

				} else if time-n.Time < -MaxRange && release(k) {
					queueScore(n)
					stagedNotes[k] = fetchNextNote(n.idx, notes)
				}

			default:
				panic("not reach")
			}
		}
		if n := copy(lastPressed, pressed); n != keymode {
			panic("copy failed")
		}
	}
}

func score(n PlayingNote) (score float64) {
	switch ScoreMode {
	case Legacy:
		// 롱노트 처리가 까다로움
		score = baseScore(n)+BonusScore(n)
	case New:
	}
	stagedNotes[n.Key] = fetchNextNote(n.idx, notes)
	return
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

func fetchNextNote(start int, notes []Note) int {
	key := notes[start].Key
	for i, n := range notes[start+1:] {
		if n.Key == key {
			return start + 1 + i
		}
	}
	return -1
}

func UnitScore(c int) float64 { return MaxScore / float64(c) }

// 굿나면 일단 굿 자체에서 절반 이상 까이고 시작
// 1굿 자체는 큰 영향 없음
// 단 미스의 경우, 미스 뒤 굿이라도 나면 25%.
// closure
func BonusScore(j Judgement) {
	var lastBonus float64
	bonus = lastBonus + j.Bonus - j.Punishment
	j.BonusValue * math.Sqrt(bonus) / 320
}
