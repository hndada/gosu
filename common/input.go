package common

import "github.com/hndada/gosu/engine/kb"

type PlayKeyEvent struct {
	Time    int64
	Pressed bool
	Key     int // Key layout index
}

func ToPlayKeyEvent(layout []kb.Code, e kb.KeyEvent) PlayKeyEvent {
	e2 := PlayKeyEvent{
		Time:    e.Time,
		Pressed: e.Pressed,
		Key:     -1,
	}
	for k, v := range layout {
		if v == e.KeyCode {
			e2.Key = k
		}
	}
	return e2
}

type KeyActionState int

const (
	Idle KeyActionState = iota
	Press
	Release
	Hold
)

// Actions are realized with 2 snapshots
func KeyAction(last, now bool) KeyActionState {
	switch {
	case !last && !now:
		return Idle
	case !last && now:
		return Press
	case last && !now:
		return Release
	case last && now:
		return Hold
	default:
		panic("not reach")
	}
}
