package mania

// go built-in inter process communcation
// https://pkg.go.dev/net/rpc

// data structure
// hook process side: list of key event log
// {time: 1432, key: 2, pressed: true}
// {time: 1501, key: 2, pressed: false}
// {time: 1592, key: 0, pressed: true}

// 1. 업데이트 때마다 hook 프로세스에서 마지막 index 이후 최신 log 불러오기
// 간단하게 생각해서 text로 취급해보자
type keyEvent struct {
	time    int64
	key     int
	pressed bool
}

// 만약 ESC 눌러서 일시정지 되었다면? (복잡하니까 일시정지는 나중에 생각: paused==false)
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

// func (s *Scene) processInput() {
// 	for {
// 		select {
// 		case <-time.After(time.Duration(s.endTime+2) * time.Millisecond):
// 			_ = s.Close()
// 		case k := <-s.kbChan.Chan:
// 			t := s.kbChan.Time().Milliseconds()
// 			for key, keyCode := range s.layout {
// 				if keyCode == k.VKCode {
// 					var e keyEvent
// 					e.time = t
// 					e.key = key
// 					// e.pressed is false by default
// 					if k.Message == types.WM_KEYDOWN {
// 						e.pressed = true
// 					}
// 					s.processScore(e)
// 				}
// 			}
// 			if k.VKCode == types.VK_ESCAPE {
// 				_ = s.Close()
// 			}
// 			// default: // skip
// 		}
// 	}
// }

// func newKeyEvent(layout []types.VKCode, e mode.KeyboardEvent) (keyEvent, bool) {
// 	for key, keyCode := range layout {
// 		if keyCode == e.KeyCode {
// 			var e2 keyEvent
// 			e2.time = e.Time
// 			e2.key = key
// 			// e2.pressed is false by default
// 			if e.State == mode.KeyStateDown {
// 				e2.pressed = true
// 			}
// 			return e2, true
// 		}
// 	}
// 	return keyEvent{}, false
// }
//
// todo: k.Time, 마지막 메시지에서 지난 시간인 듯
