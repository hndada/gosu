package mania

import "github.com/hndada/gosu/engine/kb"

type keyEvent struct {
	// kb.KeyEvent // TODO: 왜 embed가 안되지
	Time    int64
	KeyCode kb.Code
	Pressed bool
	Key     int
}

// 업데이트 때마다 마지막 index 이후 최신 log 불러오기
// TODO: ESC 눌러서 일시 정지되었을 때
// type keyEvent struct {
// 	time    int64
// 	key     int
// 	pressed bool
// }

// func FetchKeyEvents(paused *bool, reply *[]keyEvent) error {
// 	var keyEvents []keyEvent // temp
// 	var current int
// 	var count int
// 	if *paused {
// 		*reply = make([]keyEvent, 0)
// 	} else {
// 		r := make([]keyEvent, len(keyEvents)-current)
// 		for i, e := range keyEvents[current:] {
// 			r[i] = e
// 			count++
// 		}
// 		*reply = r
// 		current += count
// 	}
// 	return nil
// }
