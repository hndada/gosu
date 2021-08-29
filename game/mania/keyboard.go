package mania

// 업데이트 때마다 마지막 index 이후 최신 log 불러오기
type keyEvent struct {
	time    int64
	key     int
	pressed bool
}

// todo: ESC 눌러서 일시 정지되었을 때
func FetchKeyEvents(paused *bool, reply *[]keyEvent) error {
	var keyEvents []keyEvent // temp
	var current int
	var count int
	if *paused {
		*reply = make([]keyEvent, 0)
	} else {
		r := make([]keyEvent, len(keyEvents)-current)
		for i, e := range keyEvents[current:] {
			r[i] = e
			count++
		}
		*reply = r
		current += count
	}
	return nil
}
