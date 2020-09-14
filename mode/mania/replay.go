package mania

type Replay struct {
	// ID int64
	Score
	// 리플레이 데이터
}

// 즉석에서 계산할 것:
// hp graph
// hit error deviation

// todo: 리플레이 -> 키보드 인풋처럼
// 리플레이 구조: 마지막 status 시간, 레이아웃 키state

// todo: rg-parser로
func KeysPressed(x, keymode int) []bool {
	pressed := make([]bool, keymode)
	mask := 1
	for i := 0; i < keymode; i++ {
		pressed[i] = x&mask != 0
		mask = mask << 1 // mask *= 2
	}
	return pressed
}
