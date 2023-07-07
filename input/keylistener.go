package input

import "time"

// KeyInputListener is a listener for key input.
// Can output a replay data.
type KeyInputListener struct {
	KeySettings []Key
	vkcodes     []uint32
	PollingRate time.Duration
	Listen      func() KeyPressedLog
	// StartTime   time.Time

	Logs   []KeyPressedLog
	Index  int
	Paused bool
}

type KeyPressedLog struct {
	Time        time.Time
	PressedList []bool
}

type KeyAction int

const (
	Idle KeyAction = iota
	Hit
	Release
	Hold
)

func keyAction(last, current bool) KeyAction {
	switch {
	case !last && !current:
		return Idle
	case !last && current:
		return Hit
	case last && !current:
		return Release
	case last && current:
		return Hold
	default:
		panic("not reach")
	}
}

type KeyActionLog struct {
	Time   time.Time
	Action []KeyAction
}

func (kl *KeyInputListener) Poll() {
	if kl.Paused {
		return
	}
	log := kl.Listen()
	if !kl.isLogSame(log) {
		kl.Logs = append(kl.Logs, log)
	}
}
func (kl KeyInputListener) isLogSame(log KeyPressedLog) bool {
	lastPressed := kl.Logs[len(kl.Logs)-1].PressedList
	for i, p := range log.PressedList {
		if p != lastPressed[i] {
			return false
		}
	}
	return true
}

func (kl *KeyInputListener) Fetch() ([]KeyPressedLog, []KeyActionLog) {
	if len(kl.Logs) == 0 {
		return nil, nil
	}

	// pressedLogs
	rawLogs := kl.Logs[kl.Index:]
	pressedLogs := make([]KeyPressedLog, 0, 10)
	// now := time.Now().UnixNano()/int64(time.Millisecond) + 1
	// now := time.Now().Add(50 * time.Microsecond)
	for i, log := range rawLogs {
		start := log.Time
		var end time.Time
		if i < len(rawLogs)-1 {
			end = rawLogs[i+1].Time
		} else { // last one
			end = time.Now()
		}
		for t := start; t.Before(end); t = t.Add(kl.PollingRate) {
			pressedLogs = append(pressedLogs, KeyPressedLog{t, log.PressedList})
		}
	}

	// actionLogs
	actionLogs := make([]KeyActionLog, 0, len(pressedLogs))
	lastPressedList := make([]bool, len(kl.KeySettings))
	if kl.Index > 0 {
		lastPressedList = kl.Logs[kl.Index-1].PressedList
	}
	for _, log := range pressedLogs {
		actions := make([]KeyAction, len(log.PressedList))
		for k, p := range log.PressedList {
			actions[k] = keyAction(lastPressedList[k], p)
		}
		lastPressedList = log.PressedList
		actionLogs = append(actionLogs, KeyActionLog{log.Time, actions})
	}

	// Update Index
	kl.Index = len(kl.Logs)
	return pressedLogs, actionLogs
}
