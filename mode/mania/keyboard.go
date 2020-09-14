package mania

import (
	"github.com/moutend/go-hook/pkg/types"
	"time"
)

type keyEvent struct {
	time    int64
	key     int
	pressed bool
}

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
func (s *Scene) processInput() {
	for {
		select {
		case <-time.After(time.Duration(s.endTime+2) * time.Millisecond):
			_ = s.Close()
		case k := <-s.kbChan.Chan:
			t := s.kbChan.Time().Milliseconds()
			for key, keyCode := range s.layout {
				if keyCode == k.VKCode {
					var e keyEvent
					e.time = t
					e.key = key
					// e.pressed is false by default
					if k.Message == types.WM_KEYDOWN {
						e.pressed = true
					}
					s.processScore(e)
				}
			}
			if k.VKCode == types.VK_ESCAPE {
				_ = s.Close()
			}
		// default: // skip
		}
	}
}
