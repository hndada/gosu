package input

import "time"

// KeyListener is a listener for key input.
// Can output a replay data.
type KeyListener struct {
	KeySettings []Key
	vkcodes     []uint32 // used for windows
	PollingRate time.Duration
	Listen      func() KeyPressedLog
	StartTime   time.Time // It is used for replaying.

	PressedLogs []KeyPressedLog
	index       int
	// ActionLogs  []KeyActionLog
	// lastPressedList []bool // last pressed list since latest fetch

	Paused bool
}

type KeyPressedLog struct {
	Time        time.Time
	PressedList []bool
}

func (kl *KeyListener) Poll() {
	if kl.Paused {
		return
	}
	plog := kl.Listen() // PressedLog
	if kl.isLogSame(plog) {
		return
	}
	kl.PressedLogs = append(kl.PressedLogs, plog)
}
func (kl KeyListener) isLogSame(log KeyPressedLog) bool {
	lps := kl.lastPressed()
	for i, p := range log.PressedList {
		lp := lps[i]
		if p != lp {
			return false
		}
	}
	return true
}
func (kl KeyListener) lastPressed() []bool {
	return kl.PressedLogs[len(kl.PressedLogs)-1].PressedList
}

// func (kl *KeyListener) Fetch() []KeyActionLog {
// 	if len(kl.PressedLogs) == 0 {
// 		return nil
// 	}
// 	pressedLogs := kl.fetchPressedLogs()

// 	// Update index
// 	kl.index = len(kl.PressedLogs)
// 	return pressedLogs, actionLogs
// }
// func (kl KeyListener) fetchPressedLogs() []KeyPressedLog {
// 	rawLogs := kl.PressedLogs[kl.index:]
// 	pressedLogs := make([]KeyPressedLog, 0, 10)
// 	// now := time.Now().UnixNano()/int64(time.Millisecond) + 1
// 	// now := time.Now().Add(50 * time.Microsecond)
// 	for i, log := range rawLogs {
// 		start := log.Time
// 		var end time.Time
// 		if i < len(rawLogs)-1 {
// 			end = rawLogs[i+1].Time
// 		} else { // last one
// 			end = time.Now()
// 		}
// 		for t := start; t.Before(end); t = t.Add(kl.PollingRate) {
// 			pressedLogs = append(pressedLogs, KeyPressedLog{t, log.PressedList})
// 		}
// 	}

// 	// Append action log.
// 	alogs := make([]KeyActionLog, 0, len(pressedLogs))
// 	actions := make([]KeyActionType, len(plog.PressedList))
// 	lps := kl.lastPressed()
// 	for k, p := range plog.PressedList {
// 		actions[k] = KeyAction(lps[k], p)
// 	}
// 	alog := KeyActionLog{plog.Time, actions} // ActionLog
// 	alogs = append(alogs, alog)

// 	return pressedLogs
// }
// func (kl KeyListener) fetchActionLogs() []KeyActionLog {
// 	actionLogs := make([]KeyActionLog, 0, len(pressedLogs))
// 	lastPressedList := make([]bool, len(kl.KeySettings))
// 	if kl.index > 0 {
// 		lastPressedList = kl.PressedLogs[kl.index-1].PressedList
// 	}
// 	for _, log := range pressedLogs {
// 		actions := make([]KeyAction, len(log.PressedList))
// 		for k, p := range log.PressedList {
// 			actions[k] = keyAction(lastPressedList[k], p)
// 		}
// 		lastPressedList = log.PressedList
// 		actionLogs = append(actionLogs, KeyActionLog{log.Time, actions})
// 	}
// 	return actionLogs
// }
