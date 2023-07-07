package input

import "time"

// KeyInputListener is a listener for key input.
// Can output a replay data.
type KeyInputListener struct {
	KeySettings []Key
	vkcodes     []uint32
	PollingRate time.Duration
	Listen      func() KeyInputLog
	// StartTime   time.Time

	Logs   []KeyInputLog
	Index  int
	Paused bool
}

type KeyInputLog struct {
	Time    time.Time
	Pressed []bool
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
func (kl KeyInputListener) isLogSame(log KeyInputLog) bool {
	lastPressed := kl.Logs[len(kl.Logs)-1].Pressed
	for i, p := range log.Pressed {
		if p != lastPressed[i] {
			return false
		}
	}
	return true
}

func (kl *KeyInputListener) Fetch() []KeyInputLog {
	logs := kl.Logs[kl.Index:]
	kl.Index = len(kl.Logs)
	return logs
}
